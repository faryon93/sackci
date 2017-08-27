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
    "net/http"

    "github.com/gorilla/mux"

    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/sse"
    "github.com/faryon93/sackci/ctx"
    "github.com/faryon93/sackci/rest"
)


// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    HTTP_API_BASE = "/api/v1"
)


// --------------------------------------------------------------------------------------
//  routes
// --------------------------------------------------------------------------------------

//go:generate esc -prefix=assets -o assets.go assets
func routes(router *mux.Router) {
    api := router.PathPrefix(HTTP_API_BASE).Subrouter()
    api.NotFoundHandler = http.HandlerFunc(NotFound)
    rest.Fs = FS(false)

    // register classic REST endpoints
    api.Methods("GET").Path("/project").HandlerFunc(rest.ProjectList)
    api.Methods("GET").Path("/project/{id}/badge").HandlerFunc(rest.ProjectBadge)
    api.Methods("GET").Path("/project/{id}/trigger").HandlerFunc(rest.ProjectTrigger)
    api.Methods("GET").Path("/project/{id}/build/latest").HandlerFunc(rest.ProjectLatestBuild)

    // register REST endpoints
    rest.QueryOne(api, "/project/{Id}", ctx.Conf.Projects)
    rest.QueryAll(api, "/project/{Project}/env", model.Env{})
    rest.QueryAll(api, "/project/{Project}/history", model.Build{}, rest.QUERY_REVERSE)
    rest.QueryOne(api, "/project/{Project}/build/{Num}", model.Build{})
    rest.DeleteAll(api, "/project/{Project}/history", model.Build{}, rest.BuildPurge)

    // register SSE endpoints
    sse.Register(api, "/feed", ctx.Feed)

    // register static assets
    router.PathPrefix("/").Handler(PrettyUrl(FS(false)))
}


// --------------------------------------------------------------------------------------
//  common handlers
// --------------------------------------------------------------------------------------

func NotFound(w http.ResponseWriter, r *http.Request) {
    http.Error(w, "not found", http.StatusNotFound)
}