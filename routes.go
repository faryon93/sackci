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
//  types
// --------------------------------------------------------------------------------------

// HTTP Handler function.
type HttpFn func(http.ResponseWriter, *http.Request)


// --------------------------------------------------------------------------------------
//  routes
// --------------------------------------------------------------------------------------

func routes(router *mux.Router) {
    api := router.PathPrefix(HTTP_API_BASE).Subrouter()
    api.NotFoundHandler = http.HandlerFunc(NotFound)
    rest.Fs = FS(false)

    // register classic REST endpoints
    Get(api,"/project", rest.ProjectList)
    Get(api,"/project/{Project}/badge", rest.ProjectBadge)
    Get(api,"/project/{Project}/trigger", rest.ProjectTrigger)
    Get(api,"/project/{Project}/build/latest", rest.ProjectLatestBuild)
    Get(api,"/project/{Project}/build/{Num}/log", rest.BuildRawLog)
    Get(api,"/project/{Project}/build/{Num}/log/{stage}", rest.BuildStageLog)
    Get(api,"/project/{Project}/build/{Num}/artifacts.tar.gz", rest.BuildArtifacts)

    // the authentication endpoints should only be available when
    // authentication is enabled in the config file
    if ctx.Conf.IsAuthEnabled() {
        Post(api, "/login", rest.Login)
        Get(api, "/logout", rest.Logout)
    }

    // register model-based REST endpoints
    rest.QueryOne(api, "/project/{Id:[0-9]+}", ctx.Conf.Projects)
    rest.UpdateOne(api, "/project/{Id:[0-9]+}", ctx.Conf.Projects, rest.ProjectUpdate)
    rest.QueryAll(api, "/project/{Project:[0-9]+}/history", model.Build{}, rest.QUERY_REVERSE)
    rest.QueryOne(api, "/project/{Project:[0-9]+}/build/{Num}", model.Build{})
    rest.DeleteAll(api, "/project/{Project:[0-9]+}/history", model.Build{}, rest.BuildPurge)

    // register SSE endpoints
    sse.Register(api, "/feed", ctx.Feed)

    // register static assets
    router.PathPrefix("/").Handler(PrettyUrl(FS(false)))
}


// --------------------------------------------------------------------------------------
//  helper functions
// --------------------------------------------------------------------------------------

// Registers a handler function for the GET Method on the given path.
func Get(router *mux.Router, path string, fn HttpFn) (*mux.Route) {
    return router.Methods(http.MethodGet).Path(path).HandlerFunc(fn)
}

// Registers a handler function for the POST Method on the given path.
func Post(router *mux.Router, path string, fn HttpFn) (*mux.Route) {
    return router.Methods(http.MethodPost).Path(path).HandlerFunc(fn)
}
