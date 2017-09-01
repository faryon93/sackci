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
    "sync"

    "github.com/asdine/storm"

    "github.com/faryon93/sackci/log"
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
    Id int `yaml:"-" json:"id"`
    Name string `yaml:"name,omitempty" json:"name"`
    Scm string `yaml:"scm,omitempty" json:"scm"`
    Repository string `yaml:"repo,omitempty" json:"repository"`
    Branch string `yaml:"branch,omitempty" json:"branch"`
    Trigger string `yaml:"trigger,omitempty" json:"trigger"`
    Interval int `yaml:"interval,omitempty" json:"interval"`
    Env map[string]string `yaml:"env,omitempty" json:"env"`
    CommitUrl string `yaml:"commit_url,omitempty" json:"commit_url"`
    PrivateKey string `yaml:"key"`

    // runtime variables
    ExecutionLock sync.Mutex `json:"-" yaml:"-"`
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Creates a new Build from this project.
func (p *Project) NewBuild() (*Build) {
    // in order to assign the next build number
    // we need to get the latest build for the project
    build, err := p.GetLastBuild()
    if err != nil {
        log.Error("project", "failed to get last build:", err.Error())
    }

    // assign the next build number
    // or start with 1 if no build exists
    num := 1
    if build != nil {
        num = build.Num + 1
    }

    return &Build{
        Project: uint64(p.Id),
        Num: num,
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

// Checks the integrty of a project and corrects if necessary
func (p *Project) CheckIntegrity() {
    // query all builds
    var builds []Build
    err := Get().Find("Project", p.Id, &builds)
    if err == storm.ErrNotFound {
        return
    } else if err != nil {
        log.Error("project", "integrity check failed:", err.Error())
        return
    }

    for _, build := range builds {
        if build.Status == BUILD_RUNNING {
            log.Info("project", "found build", build.Num, " of project \"" + p.Name +
                "\" which is still \"" + build.Status + "\". Canceling the build.")

            // cancel the build
            build.Status = BUILD_FAILED
            build.Save()
        }
    }
}