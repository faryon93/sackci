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
    log "github.com/sirupsen/logrus"

    "github.com/faryon93/sackci/ctx"
    "github.com/faryon93/sackci/agent"
    "github.com/faryon93/sackci/model"
    "errors"
)


// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    TOKEN_QUERY_NAME = "token"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type triggerResponse struct {
    Success bool `json:"success"`
    BuildId int `json:"build_id"`
}


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func ProjectTrigger(w http.ResponseWriter, r *http.Request) {
    // if the pipeline is executed, the project has to be
    // unlocked when the pipeline finished
    ignoreDeferedUnlock := false

    // parse the id of the project
    id, err := strconv.Atoi(mux.Vars(r)["Project"])
    if err != nil {
        http.Error(w, "invalid project id", http.StatusNotAcceptable)
        return
    }

    // search for the project in the configuration
    project := ctx.Conf.GetProject(id)
    if project == nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }

    // if trigger tokens are definied we need to
    // check if the user supplied token is valid
    err = isTriggerAllowed(project, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    // try to lock the project -> if fails, a build is already running
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
        log.Errorln("failed to trigger build:", err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    pipeline.SetProject(project)
    pipeline.SetArtifactsDir(ctx.Conf.GetArtifactsDir())

    // construct the build object for saving in the database
    build := project.NewBuild()
    build.Time = pipeline.StartTime
    build.Node = pipeline.Agent.Name
    build.Stages = []model.Stage{
        {Name: agent.STAGE_SCM_NAME, Status: model.STAGE_IGNORED, Log: []string{}},
    }
    err = build.Save()
    if err != nil {
        log.Errorln("failed to save build:", err.Error())
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

    Jsonify(w, triggerResponse{true, build.Num})
}


// --------------------------------------------------------------------------------------
//  private functions
// --------------------------------------------------------------------------------------

// Returns nil when the calling request is allowed to perform a build trigger.
// Otherwise an error with the reason is returned.
func isTriggerAllowed(project *model.Project, r *http.Request) (error) {
    // trigger is allowed at anytime when authentication is turned off
    if !ctx.Conf.IsAuthEnabled() {
        return nil
    }

    queryToken := r.URL.Query().Get(TOKEN_QUERY_NAME)
    accessGranted := false
    if len(queryToken) > 0 {
        // check if the supplied token is in the token list
        // in order to grad access
        for _, token := range project.TriggerTokens {
            if queryToken == token {
                accessGranted = true
                break
            }
        }

        // the provided token is invalid -> tell the caller
        if !accessGranted {
            return errors.New("invalid trigger token")
        }
    }

    // the token based authentication was not successfull
    // next try the session based authentication
    if !accessGranted {
        _, err := ctx.Sessions.ValiadeRequest(r)
        if err != nil {
            return err
        }
    }

    return nil
}