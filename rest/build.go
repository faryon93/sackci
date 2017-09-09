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
    "path/filepath"
    "fmt"
    "os"
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

// Gets the whole raw log of a build.
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

// Gets the raw log of a specific stage of a build.
func BuildStageLog(w http.ResponseWriter, r *http.Request) {
    stage, err := strconv.Atoi(mux.Vars(r)["stage"])
    if err != nil {
        http.Error(w, "invalid stage", http.StatusNotAcceptable)
        return
    }

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

    // check stage boundries
    if stage < 0 ||stage >= len(build.Stages) {
        http.Error(w, "invalid stage", http.StatusNotAcceptable)
        return
    }

    w.Header().Set("Content-Type", CONTENT_TYPE_TEXT)
    w.Write([]byte(build.Stages[stage].RawLog()))
}

// Gets the tar.gz archive of the build artifacts.
func BuildArtifacts(w http.ResponseWriter, r *http.Request) {
    project := mux.Vars(r)["Project"]
    build := mux.Vars(r)["Num"]

    if val , err := strconv.Atoi(project); val <= 0 || err != nil {
        http.Error(w, "invalid project", http.StatusNotAcceptable)
        return
    }

    if val, err := strconv.Atoi(build); val <= 0 || err != nil {
        http.Error(w, "invalid build number", http.StatusNotAcceptable)
        return
    }

    // check if the artifact exists in filesystem
    artifactFile := filepath.Join(ctx.Conf.GetArtifactsDir(), project, build + ".tar.gz")
    if _, err := os.Stat(artifactFile); err != nil {
        http.Error(w, "artifact not found", http.StatusNotFound)
        return
    }

    // set headers to start a download
    filename := fmt.Sprintf("artifact-%s-%s.tar.gz", project, build)
    w.Header().Set("Content-Type", CONTENT_TYPE_STREAM)
    w.Header().Set("Content-Disposition", "attachment; filename=" + filename)
    http.ServeFile(w, r, artifactFile)
}

// Callback when build history purging was successfull.
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
