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

    "github.com/asdine/storm"
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
        build, err := getLastBuild(id)
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

        // no build has executed yet -> assign the correct status
        } else {
            item.Status = model.BUILD_STATUS_WAITING
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


// --------------------------------------------------------------------------------------
//  private functions
// --------------------------------------------------------------------------------------

func getLastBuild(project int) (*model.Build, error) {
    bolt := model.Get()

    // fetch the last inserted build for the project
    var builds []model.Build
    err := bolt.Find("Project", project, &builds, storm.Limit(1), storm.Reverse())
    if err == storm.ErrNotFound {
        return nil, nil

    // an internal error occoured
    } else if err != nil {
        return nil, err
    }

    // the query returned a resulst -> return it
    if len(builds) > 0 {
        return &builds[0], nil

    // no build was found
    } else {
        return nil, nil
    }
}