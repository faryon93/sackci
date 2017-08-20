package model
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

// ----------------------------------------------------------------------------------
//  imports
// ----------------------------------------------------------------------------------

import (
    "github.com/asdine/storm"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    TRIGGER_MANUAL = "manual"
    TRIGGER_POLL   = "poll"
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type Project struct {
    Id int
    Name string `yaml:"name"`
    Scm string `yaml:"scm"`
    Repository string `yaml:"repo"`
    Branch string `yaml:"branch"`
    Trigger string `yaml:"trigger"`
    Interval int `yaml:"interval"`
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Creates a new Build from this project.
func (p *Project) NewBuild() (*Build) {
    return &Build{
        Project: uint64(p.Id),
        Num: 4,
        Status: BUILD_RUNNING,
        Commit: Commit{
            Message: "unknown",
            Author: "unknown",
            Ref: "unknown",
        },
    }
}

// Gets the latest build of this project.
func (p *Project) GetLastBuild() (*Build, error) {
    // fetch the last inserted build for the project
    var builds []Build
    err := Get().Find("Project", p.Id, &builds, storm.Limit(1), storm.Reverse())
    if err == storm.ErrNotFound {
        return nil, nil

    // an internal error occoured
    } else if err != nil {
        return nil, err
    }

    // no build was found
    // this case should be handled by err == strom.ErrNotFound
    if len(builds) <= 0 {
        return nil, nil
    }

    return &builds[0], nil
}