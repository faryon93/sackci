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
    "time"
)


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

func Login(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "faild to parse login data", http.StatusNotAcceptable)
        return
    }

    // validate username and password
    if r.Form.Get("username") == "test" &&
       r.Form.Get("password") == "test" {


        http.SetCookie(w, &http.Cookie{
            Path: "/", Name: "session",
            Value: "test", Secure: true,
        })
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

    // the password / username matching failed -> return an error
    } else {
        http.Redirect(w, r, "/login?error=invalid_credentials", http.StatusSeeOther)
    }
}

func Logout(w http.ResponseWriter, r *http.Request) {
    // just send an already expired and empty cookie
    // back to the client
    // when we return the Unauthorized HTTP code
    // the frontend will rediret to the login page itself
    http.SetCookie(w, &http.Cookie{
        Name: "session", Value: "", Secure: true, Path: "/", Expires: time.Now().Add(-5 * time.Minute),
    })
    http.Error(w, "logout successfull", http.StatusUnauthorized)
}
