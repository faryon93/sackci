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

    "github.com/faryon93/sackci/assets"
//    "strings"
    "strings"
)


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

// Default not found handler.
func NotFound(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "not found", http.StatusNotFound)
}

// Redirects to the configured https endpoint.
func RedirectHttps(w http.ResponseWriter, r *http.Request) {
    url := "https://" + r.Host + r.URL.String()
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Rewrites everything which is not contained in the embedded asset fs
// to the root path in order to support HTML5 in browser naivation.
func PrettyUrl(fs http.FileSystem) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // do not serve lua files from assetfs
        if strings.HasSuffix(r.URL.Path, ".lua") {
            http.Error(w, "not found", http.StatusNotFound)
            return
        }

        // the webpage root cannot be checked in fs
        // everything without a file ending should to be
        // rewritten to index page for a more convenient
        // user experiance and to prevent the server from delivering
        // wrong file content to expected filetype
        if r.URL.Path != "/" && !assets.FileExists(r.URL.String()) {
            r.URL.Path = "/"
        }

        http.FileServer(fs).ServeHTTP(w, r)
    })
}
