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

    "github.com/faryon93/sackci/ctx"
    "github.com/faryon93/sackci/assets"
)


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

// Sessions Middleware.
// Checks if a proper session token is supplied by the caller.
func CheckSession(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        url := r.URL.Path
        token, err := ctx.Sessions.ValiadeRequest(r)
        if err == nil {
            ctx.Sessions.Refresh(token)
        }

        // the whole rest api is secured by a session token
        // except the login endpoint
        // the trigger endpoint handls authentication himself
        if strings.HasPrefix(url, HTTP_API_BASE) {
            secured := url != HTTP_API_BASE + "/login" &&
                     !strings.HasSuffix(url, "/trigger")

            if secured && err != nil {
                http.Error(w, err.Error(), http.StatusUnauthorized)
                return
            }

        // all non api pages should be redirected to the login page
        // except for all static assets
        } else if !strings.HasPrefix(url, "/login") && !assets.FileExists(url) {
            if err != nil {
                http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
                return
            }
        }

        h.ServeHTTP(w, r)
    })
}
