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
    "time"
    "strconv"

    "github.com/gorilla/mux"

    "github.com/faryon93/sackci/ctx"
    "github.com/faryon93/sackci/model"
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func BuildRawLog(w http.ResponseWriter, r *http.Request) {
    // construct the query
    query, err := stormQuery(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotAcceptable)
        return
    }

    var build model.Build
    err = query.First(&build)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", CONTENT_TYPE_TEXT)

    for _, stage := range build.Stages {
        w.Write([]byte(stage.RawLog()))
    }
}

func BuildPurge(r *http.Request) {
    // parse the url parameters for the id
    fieldVal, err := strconv.Atoi(mux.Vars(r)["Project"])
    if err != nil {
        return
    }

    // publish the event, which informs about the change
    // the Build Number of 0 does not exist, so it is clear what happend
    ctx.Feed.Publish(model.EvtPipelineFinished{
        EventBase: &model.EventBase{
            Project: fieldVal,
            Build: 0,
            Timestamp: time.Now().UnixNano(),
        },
        Status: model.BUILD_WAITING,
        Duration: 0,
    })
}