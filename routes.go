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
//  imports
// --------------------------------------------------------------------------------------

import (
    "github.com/gorilla/mux"

    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/sse"
    "github.com/faryon93/sackci/ctx"
    "github.com/faryon93/sackci/rest"
)


// --------------------------------------------------------------------------------------
//  routes
// --------------------------------------------------------------------------------------

func routes(router *mux.Router) {
    // register REST endpoints
    rest.All(router, "/project", model.Project{})
    rest.One(router, "/project/{id}", model.Project{})
    rest.One(router, "/build/{id}", model.Build{})
    rest.QueryAll(router, "/project/{project}/env","Project", model.Env{})
    rest.QueryAll(router, "/project/{project}/history", "Project", model.Build{})

    // register SSE endpoints
    sse.Register(router, "/feed", ctx.Feed)
}