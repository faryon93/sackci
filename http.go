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
    "time"
    "context"

    log "github.com/sirupsen/logrus"

    "github.com/faryon93/sackci/config"
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
        srv = &http.Server{Addr: conf.HttpListen, Handler: getHandler(conf, mux)}
    }

    go func() {
        log.Infoln("http server is listening on http://" + conf.HttpListen)

        // serve the http connection as configured
        err := srv.ListenAndServe()
        if err != nil && err != http.ErrServerClosed {
            log.Errorln("failed to serv http:", err.Error())
            return
        }

        log.Infoln("http server is now closed")
    }()

    return srv
}

// Starts an TLS encrypted http endpoint.
func SetupHttpsEndpoint(conf *config.Config, mux http.Handler) (*http.Server) {
    srv := &http.Server{Addr: conf.HttpsListen, Handler: getHandler(conf, mux)}
    go func() {
        log.Infoln("https server is listening on https://" + conf.HttpsListen)

        // serve the http connection as configured
        err := srv.ListenAndServeTLS(conf.HttpsCert, conf.HttpsKey)
        if err != nil && err != http.ErrServerClosed {
            log.Errorln("failed to serv https:", err.Error())
            return
        }

        log.Infoln("https server is now closed")
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

func getHandler(conf *config.Config, mux http.Handler) (http.Handler) {
    // only apply session check if authentication is enabled
    handler := mux
    if conf.IsAuthEnabled() {
        handler = CheckSession(mux)
    }

    return handler
}