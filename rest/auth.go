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
    "github.com/faryon93/sackci/ctx"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    // Time the session cookie should be valid when remember
    // me option is used during login.
    SESSION_TIMEOUT_REMEBERME = 7 * 24 * time.Hour
    SESSION_TIMEOUT_NORMAL    = 15 * time.Minute

    // Name of the session cookie.
    SESSION_COOKIE = "token"
)


// ----------------------------------------------------------------------------------
//  global variables
// ----------------------------------------------------------------------------------

var (
    Sessions = NewSessionStore()
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

// Try to login the user with provided credentials.
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
        // if the rember me option is not selected we have
        // a defined session timeout. The session will be refereshed
        // on each request don by the user.
        timeout := SESSION_TIMEOUT_NORMAL
        if login.Remeber {
            timeout = SESSION_TIMEOUT_REMEBERME
        }

        // create the new session
        token, err := Sessions.Create(timeout)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // create the session cookie
        cookie := http.Cookie{
            Path:     "/", Name: SESSION_COOKIE,
            Value:    token,
            HttpOnly: true,
            Secure:   ctx.Conf.IsHttpsEnabled(),
        }

        // set an expiration time on the cookie
        // if the remeber me option is used
        if login.Remeber {
            cookie.Expires = time.Now().Add(timeout)
        }

        // send the cookie back to client in the http response
        http.SetCookie(w, &cookie)

    // the password / username matching failed -> return an error
    } else {
        http.Error(w, "Invalid username or password", http.StatusNotAcceptable)
        return
    }
}

// Logout the user.
func Logout(w http.ResponseWriter, r *http.Request) {
    // trying to logout without a session cookie should
    // result in the same behaviour as with a valid session
    token, err := GetCookie(r, SESSION_COOKIE)
    if err == nil {
        Sessions.Delete(token)
    }

    // TODO: how to handle SSE disconnect?

    // Just send an already expired and empty cookie back to the client.
    // When we return the Unauthorized HTTP code the frontend
    // will rediret to the login page automatically
    http.SetCookie(w, InvalidateCookie("/", SESSION_COOKIE))
    http.Error(w, "logout successfull", http.StatusUnauthorized)
}
