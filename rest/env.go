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
    "github.com/faryon93/sackci/model"
)


// --------------------------------------------------------------------------------------
//  http handlers
// --------------------------------------------------------------------------------------

func EnvGet(w http.ResponseWriter, r *http.Request) {
    project, err :=  strconv.Atoi(mux.Vars(r)["project"])
    if err != nil {
        http.Error(w, "invalid project id", http.StatusNotAcceptable)
        return
    }

    // ignore the errors
    env,_ := model.GetProjectEnv(uint64(project))
    Jsonify(w, env)
}
