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

// --------------------------------------------------------------------------------------
// imports
// --------------------------------------------------------------------------------------

import (
    "net/http"
    "github.com/gorilla/mux"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

// Composition of multiple Routes.
type Routes []Route

// A simple route which can be registered.
type Route struct {
    Method string
    Path string
    Fn func(w http.ResponseWriter, r *http.Request)
}


// --------------------------------------------------------------------------------------
// public members
// --------------------------------------------------------------------------------------

func (r Routes) Setup(router *mux.Router) {
    for _, route := range r {
        router.Methods(route.Method).Path(route.Path).HandlerFunc(route.Fn)
    }
}
