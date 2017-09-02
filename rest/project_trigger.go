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
    "strconv"

    "github.com/gorilla/mux"

    "github.com/faryon93/sackci/ctx"
    "github.com/faryon93/sackci/agent"
    "github.com/faryon93/sackci/log"
    "github.com/faryon93/sackci/model"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type triggerResponse struct {
    Success bool `json:"success"`
    BuildId int `json:"build_id"`
    Message string `json:"message"`
}


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func ProjectTrigger(w http.ResponseWriter, r *http.Request) {
    // if the pipeline is executed, the project has to be
    // unlocked when the pipeline finished
    ignoreDeferedUnlock := false

    // parse the id of the project
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "invalid project id", http.StatusNotAcceptable)
        return
    }

    // return the runtime object
    project := ctx.Conf.GetProject(id)
    if project == nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }

    // try to lock the project -> if fails, a build
    // is already running
    if err := project.Lock(); err != nil {
        http.Error(w, "Build already running", http.StatusConflict)
        return
    }
    defer func() {
        if !ignoreDeferedUnlock {
            project.Unlock()
        }
    }()

    // create the pipeline on the build agent
    pipeline, err := agent.CreatePipeline()
    if err != nil {
        log.Error("project", "failed to trigger build:", err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    pipeline.SetProject(project)

    // construct the build object for saving in the database
    build := project.NewBuild()
    build.Time = pipeline.StartTime
    build.Node = pipeline.Agent.Name
    build.Stages = []model.Stage{
        {Name: agent.STAGE_SCM_NAME, Status: model.STAGE_IGNORED, Log: []string{}},
    }
    err = build.Save()
    if err != nil {
        log.Error("project", "failed to save build:", err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    pipeline.SetBuild(build)

    // redirect and transform all events for the eventstream
    go func() {
        for event := range pipeline.Events {
            build.Publish(event)
            ctx.Feed.Publish(event)
        }
    }()

    // asynchrounsly execute the proejct on the provisioned pipeline
    ignoreDeferedUnlock = true
    go func() {
        pipeline.Execute()
        pipeline.Destroy()

        // release to project -> another build can now be executed
        project.Unlock()
    }()

    Jsonify(w, triggerResponse{
        Success: true,
        BuildId: build.Num,
    })
}
