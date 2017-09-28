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

// ----------------------------------------------------------------------------------
//  imports
// ----------------------------------------------------------------------------------

import (
    "net/http"
    "io/ioutil"
    "encoding/json"
    "time"
)

// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type loginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Remeber bool `json:"remember"`
}


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

func Login(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // parse the login information
    var login loginRequest
    err = json.Unmarshal(body, &login)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotAcceptable)
        return
    }

    // validate username and password
    if login.Username == "test" && login.Password == "test" {
        // create the session cookie and send it
        // in the http response
        http.SetCookie(w, &http.Cookie{
            Path: "/", Name: "session",
            Value: "test", HttpOnly: true,
        })

    // the password / username matching failed -> return an error
    } else {
        http.Error(w, "Invalid username or password", http.StatusNotAcceptable)
        return
    }
}

func Logout(w http.ResponseWriter, r *http.Request) {
    // Just send an already expired and empty cookie back to the client.
    // When we return the Unauthorized HTTP code the frontend
    // will rediret to the login page automatically
    http.SetCookie(w, &http.Cookie{
        Path: "/", Name: "session",
        Value: "",
        HttpOnly: true,
        Expires: time.Now().Add(-5 * time.Minute),
        MaxAge: -1,
    })
    http.Error(w, "logout successfull", http.StatusUnauthorized)
}
