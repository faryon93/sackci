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
//  types
// --------------------------------------------------------------------------------------

type projectListItem struct {
    Id       int             `json:"id"`
    Name     string          `json:"name"`
    Status   string          `json:"status"`
    BuildId  int             `json:"build"`
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
    // TODO: some smarter way to generate the project ids
    list := make([]projectListItem, len(projects))
    for index, project := range projects {
        // the project ids start with id = 1
        id := index + 1

        // find the latest build for the project
        build, err := project.GetLastBuild()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // construct the list entry
        item := projectListItem{
            Id: id,
            Name: project.Name,
        }

        // a build is present -> display it
        if build != nil {
            item.Status = build.Status
            item.BuildId = int(build.Id)
            item.Time = build.Time
            item.Duration = build.Duration
            item.BuildNum = build.Num

        // no build has executed yet -> assign the correct status
        } else {
            item.Status = model.BUILD_WAITING
        }

        // add the new item to the list
        list[index] = item
    }

    Jsonify(w, list)
}

// Queries on project by its id.
func ProjectById(w http.ResponseWriter, r *http.Request) {
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

    Jsonify(w, projectListItem{
        Id: id,
        Name: project.Name,
    })
}

// Queries on project by its id.
func ProjectLatestBuild(w http.ResponseWriter, r *http.Request) {
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

    build, err := project.GetLastBuild()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if build == nil {
        http.Error(w, "no build yet", http.StatusNotFound)
        return
    }

    Jsonify(w, build)
}
