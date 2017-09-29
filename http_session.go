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
    "github.com/faryon93/sackci/rest"
    "errors"
)


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

// Session Middleware.
// Checks if a proper session token is supplied by the caller.
func CheckSession(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        url := r.URL.String()
        err := validateSession(r)

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


// ----------------------------------------------------------------------------------
//  private functions
// ----------------------------------------------------------------------------------

// Checks if the session token is valid.
func validateSession(r *http.Request) (error) {
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
