package pipelinefile
// dockertest
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

// ----------------------------------------------------------------------------------
//  imports
// ----------------------------------------------------------------------------------

import (
    "encoding/json"
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type Definition struct {
    Stages []Stage `json:"pipeline"`
    Artifacts string `json:"artifacts,omitempty"`
}


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

func Parse(file []byte) (*Definition, error) {
    var buildfile Definition
    return &buildfile, json.Unmarshal(file, &buildfile)
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (b *Definition) StageString() (string) {
    stages := ""
    for i, stage := range b.Stages {
        stages += stage.Name

        if i < len(b.Stages) - 1 {
            stages += ", "
        }
    }

    return stages
}