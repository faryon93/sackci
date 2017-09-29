package rest
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

// --------------------------------------------------------------------------------------
//  imports
// --------------------------------------------------------------------------------------

import (
    "net/http"
    "encoding/json"
    "io"
    "time"
    "errors"
)


// --------------------------------------------------------------------------------------
//  global variables
// --------------------------------------------------------------------------------------

var (
    Fs http.FileSystem
)


// --------------------------------------------------------------------------------------
//  imports
// --------------------------------------------------------------------------------------

const (
    CONTENT_TYPE_JSON   = "application/json"
    CONTENT_TYPE_SVG    = "image/svg+xml"
    CONTENT_TYPE_TEXT   = "text/plain"
    CONTENT_TYPE_STREAM  = "application/octet-stream"
)

// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

// Writes the JSON representation of v to the supplied http.ResposeWriter.
// If an error occours while marshalling the object the http response
// will be an internal server error.
func Jsonify(w http.ResponseWriter, v interface{}) {
    js, err := json.Marshal(v)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
    w.Write(js)
}

// Serves a file from the global assetfs.
func ServeFile(w http.ResponseWriter, path string, contentType string) {
    if Fs == nil {
        return
    }

    // proper headers
    w.Header().Set("Content-Type", contentType)

    // load the file from the asset fs
    file, err := Fs.Open(path)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    // copy all bytes of the file to the http output writer
    _, err = io.Copy(w, file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

// Applies necessary headers to disable caching.
func NoCaching(w http.ResponseWriter) {
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")
}

// Gets the value of a cookie.
func GetCookie(r *http.Request, name string) (string, error) {
    cookie, err := r.Cookie(name)
    if err != nil {
        return "", err
    }

    if cookie == nil {
        return "", errors.New("no cookie")
    }

    return cookie.Value, nil
}

// Invalidates the given cookie.
func InvalidateCookie(path, name string) (*http.Cookie) {
    return &http.Cookie{
        Path:     path,
        Name:     name,
        Value:    "",
        HttpOnly: true,
        Expires:  time.Now().Add(-5 * time.Minute),
        MaxAge:   -1,
    }
}