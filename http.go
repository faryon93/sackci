package main
// sackci
// Copyright (C) 2017 Maximilian Pachl

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// ----------------------------------------------------------------------------------
//  imports
// ----------------------------------------------------------------------------------

import (
    "net/http"
    "strings"
    "errors"
    "time"
    "context"
    "net"

    "github.com/faryon93/sackci/log"
    "github.com/faryon93/sackci/config"
    "github.com/faryon93/sackci/ctx"
    "github.com/faryon93/sackci/rest"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    LOG_TAG_HTTP = "http"
)


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

// Starts the http endpoint.
// If a TLS encrypted endpoint is configured this endpoint is just used
// to redirect automatically to the secured endpoint.
func SetupHttpEndpoint(conf *config.Config, mux http.Handler) (*http.Server) {
    var srv *http.Server

    // if https is enabled the http server is responsible for
    // redirecting unsecured request to a secure endpoint only!
    if conf.IsHttpsEnabled() {
        srv = &http.Server{
            Addr: conf.HttpListen,
            Handler: http.HandlerFunc(RedirectHttps),
        }

    // serve the normal api and frontend endpoints
    } else {
        srv = &http.Server{Addr: conf.HttpListen, Handler: CheckSession(mux)}
    }

    go func() {
        log.Info(LOG_TAG_HTTP, "http server is listening on http://" + conf.HttpListen)

        // serve the http connection as configured
        err := srv.ListenAndServe()
        if err != nil && err != http.ErrServerClosed {
            log.Error(LOG_TAG_HTTP, "failed to serv http:", err.Error())
            return
        }

        log.Info(LOG_TAG_HTTP, "http server is now closed")
    }()

    return srv
}

// Starts an TLS encrypted http endpoint.
func SetupHttpsEndpoint(conf *config.Config, mux http.Handler) (*http.Server) {
    srv := &http.Server{Addr: conf.HttpsListen, Handler: CheckSession(mux)}
    go func() {
        log.Info(LOG_TAG_HTTP, "https server is listening on https://" + conf.HttpsListen)

        // serve the http connection as configured
        err := srv.ListenAndServeTLS(conf.HttpsCert, conf.HttpsKey)
        if err != nil && err != http.ErrServerClosed {
            log.Error(LOG_TAG_HTTP, "failed to serv https:", err.Error())
            return
        }

        log.Info(LOG_TAG_HTTP, "https server is now closed")
    }()

    return srv
}

// Gracefully destroys a http endpoint with the given timeout.
func ShutdownHttp(srv *http.Server, timeout time.Duration) {
    httpCtx, _ := context.WithTimeout(context.Background(), timeout)
    srv.Shutdown(httpCtx)
}


// ----------------------------------------------------------------------------------
//  private functions
// ----------------------------------------------------------------------------------

// Redirects to the configured https endpoint.
func RedirectHttps(w http.ResponseWriter, r *http.Request) {
    host, _, err := net.SplitHostPort(r.Host)
    if err != nil {
        host = r.Host
    }

    // replace port with the configured https port
    _, httpsPort, err := net.SplitHostPort(ctx.Conf.HttpsListen)
    if err != nil {
        log.Error(LOG_TAG_HTTP,"invalid value in https_listen property")
        http.Error(w, "misconfigured server", http.StatusInternalServerError)
        return
    }
    if httpsPort != "443" {
        host = net.JoinHostPort(host, httpsPort)
    }

    url := "https://" + host + r.URL.String()
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Checks if the session token is valid.
func ValidateSession(r *http.Request) (error) {
    cookie, err := r.Cookie(rest.SESSION_COOKIE)
    if err != nil || cookie == nil {
        return errors.New("no session cookie")
    }

    if !rest.Sessions.IsValid(cookie.Value) {
        return errors.New("invalid session token")
    }

    // the session is valid -> we can refresh the token
    rest.Sessions.Refresh(cookie.Value)

    return nil
}

// Session Middleware.
// Checks if a proper session token is supplied by the caller.
func CheckSession(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        url := r.URL.String()
        err := ValidateSession(r)

        // the whole rest api is secured by a session token
        // except the login endpoint
        if strings.HasPrefix(url, HTTP_API_BASE) {
            if  url != HTTP_API_BASE + "/login" && err != nil {
                http.Error(w, err.Error(), http.StatusUnauthorized)
                return
            }

        // all non api pages should be redirected to the login page
        // except for all static assets
        } else if !strings.HasPrefix(url, "/login") && !AssetFileExists(url) {
            if err != nil {
                http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
                return
            }
        }

        h.ServeHTTP(w, r)
    })
}
