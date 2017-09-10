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
    "errors"

    "github.com/gorilla/mux"

    "github.com/faryon93/sackci/ctx"
    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/log"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type projectListItem struct {
    Id       int             `json:"id"`
    Name     string          `json:"name"`
    Status   string          `json:"status"`
    BuildNum int             `json:"build_num"`
    Time     time.Time       `json:"execution_time"`
    Duration time.Duration   `json:"duration"`
}


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

// Gets the project short list.
func ProjectList(w http.ResponseWriter, r *http.Request) {
    projects := ctx.Conf.Projects

    // construct the project list
    list := make([]projectListItem, 0)
    for _, project := range projects {
        if !project.IsValid() {
            continue
        }

        // find the latest build for the project
        build, err := project.GetLastBuild()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // construct the list entry
        item := projectListItem{
            Id: project.Id,
            Name: project.Name,
        }

        // a build is present -> display it
        if build != nil {
            item.Status = build.Status
            item.Time = build.Time
            item.Duration = build.Duration
            item.BuildNum = build.Num

        // no build has executed yet -> assign the correct status
        } else {
            item.Status = model.BUILD_WAITING
        }

        // add the new item to the list
        list = append(list, item)
    }

    Jsonify(w, list)
}

// Queries on project by its id.
func ProjectLatestBuild(w http.ResponseWriter, r *http.Request) {
    build, err := getLastBuild(w, r)
    if err != nil {
        return
    }

    if build == nil {
        http.Error(w, "no build yet", http.StatusNotFound)
        return
    }

    Jsonify(w, build)
}

// Badge (failing / passing) for the project.
func ProjectBadge(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "invalid project id", http.StatusNotAcceptable)
        return
    }

    // get the project
    project := ctx.Conf.GetProject(id)
    if project == nil {
        http.Error(w, "project not found", http.StatusNotFound)
        return
    }

    // send an error if the badge is disabled
    if !project.BadgeEnable {
        http.Error(w, "not enabled", http.StatusNotFound)
        return
    }

    build, err := project.GetLastBuild()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // the correct badge to display
    file := "/img/waiting.svg"
    if build != nil {
        file = "/img/" + build.Status + ".svg"
    }

    // serve the file and disable all browser side caching
    // its okay to include this badge from anywhere -> disable CORS
    w.Header().Set("Access-Control-Allow-Origin", "*")
    NoCaching(w)
    ServeFile(w, file, CONTENT_TYPE_SVG)
}

// Success Handler: Update Project
func ProjectUpdate(r *http.Request) {
    err := ctx.Conf.Save()
    if err != nil {
        log.Info("project", "failed to write config file:", err.Error())
        return
    }
    log.Info("project", "succesfully updated project with id: \"" + mux.Vars(r)["Id"] + "\"")
}


// --------------------------------------------------------------------------------------
//  private helpers
// --------------------------------------------------------------------------------------

func getLastBuild(w http.ResponseWriter, r *http.Request) (*model.Build, error) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "invalid project id", http.StatusNotAcceptable)
        return nil, err
    }

    // get the project
    project := ctx.Conf.GetProject(id)
    if project == nil {
        http.Error(w, "project not found", http.StatusNotFound)
        return nil, errors.New("project not found")
    }

    build, err := project.GetLastBuild()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return nil, err
    }

    return build, nil
}