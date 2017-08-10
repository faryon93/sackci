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
    "github.com/gorilla/mux"
    "github.com/faryon93/sackci/model"

    "net/http"
    "strconv"
)


// --------------------------------------------------------------------------------------
//  http handlers
// --------------------------------------------------------------------------------------

func ProjectList(w http.ResponseWriter, r *http.Request) {
    projects, err := model.QueryProjects()
    if err != nil {
        http.Error(w, "failed to fetch project list: " + err.Error(), http.StatusNotAcceptable)
        return
    }

    Jsonify(w, projects)
}

func ProjectGet(w http.ResponseWriter, r *http.Request) {
    id, err :=  strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "invalid project id", http.StatusNotAcceptable)
        return
    }

    // get the project from the database
    var project model.Project
    err = model.Get(model.PROJECTS_BUCKET, uint64(id), &project)
    if err != nil {
        http.Error(w, "failed to fetch project: " + err.Error(), http.StatusNotFound)
        return
    }

    Jsonify(w, project)
}
