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
    "errors"

    "github.com/asdine/storm"

    "github.com/faryon93/sackci/log"
    "github.com/faryon93/sackci/util"
    "strings"
    "io/ioutil"
    "net/url"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    TRIGGER_MANUAL  = "manual"
    TRIGGER_POLL    = "poll"
    KEY_PREFIX      = "-----BEGIN"
)

var (
    ErrBuildRunning = errors.New("build running")
    ErrUnknownHash = errors.New("unknown hash")
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
    BadgeEnable bool `yaml:"badge" json:"badge"`
    Hash string `yaml:"hash,omitempty" json:"-"`
    PrivateKey string `yaml:"key,omitempty" json:"key"`

    // runtime variables
    mutex sync.Mutex `json:"-" yaml:"-"`
    buildMutex sync.Mutex `json:"-" yaml:"-"`
    buildRunning bool `yaml:"-" json:"-"`
}

type ProjectMapping struct {
    Id int `storm:"id,increment"`
    Hash string `storm:"unique"`
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Creates a new Build from this project.
func (p *Project) NewBuild() (*Build) {
    p.buildMutex.Lock()
    defer p.buildMutex.Unlock()

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

// Assigns a proper ID to this project.
func (p *Project) AssignId() (error) {
    // if this project was added recently there is no
    // hash -> generate one and store mapping in database
    if util.StrEmpty(p.Hash) {
        hash := util.Hash(*p)

        // insert the mapping into the database
        mapping := ProjectMapping{Hash: hash}
        err := Get().Save(&mapping)
        if err != nil {
            log.Error("project", "failed to initialize new project \"" + p.Name + "\"")
            return err
        }

        // update the project runtime object
        p.Hash = hash
        p.Id = mapping.Id

        log.Info("project", "found new project \"" + p.Name + "\", assigning new hash", hash)
        return nil

    // there is a hash in the config file -> lookup id
    } else {
        var mapping ProjectMapping
        err := Get().One("Hash", p.Hash, &mapping)
        if err == storm.ErrNotFound {
            return ErrUnknownHash
        } else if err != nil {
            return err
        }

        p.Id = mapping.Id
    }

    return nil
}

// Checks the integrty of a project and corrects if necessary.
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
            log.Info("project", "found build", build.Num, "of project \"" + p.Name +
                "\" which is still \"" + build.Status + "\". Canceling the build.")

            // cancel the running stages
            for y := range build.Stages {
                if build.Stages[y].Status == STAGE_RUNNING {
                    build.Stages[y].Status = STAGE_FAILED
                }
            }

            // cancel the build
            build.Status = BUILD_FAILED
            build.Save()
        }
    }
}

// Marks the project as running.
func (p *Project) Lock() (error) {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    // if a build for this project is already running
    // an error should be thrown
    if p.buildRunning {
        return ErrBuildRunning
    }

    p.buildRunning = true

    return nil
}

// Marks the project as finished.
func (p *Project) Unlock() {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    p.buildRunning = false
}

// Returns true if the project is valid and should be
// made available to the public.
func (p *Project) IsValid() (bool) {
    return p.Id > 0
}

// Returns the private key which is configured
// for the project. It is ether read from the
// given file or from the property itself.
func (p *Project) GetPrivateKey() ([]byte, error) {
    if len(p.PrivateKey) <= 0 {
        return []byte{}, nil
    }

    // check if the confile file contains the
    // key directoy as plain text
    if strings.HasPrefix(strings.TrimSpace(p.PrivateKey), KEY_PREFIX) {
        return []byte(p.PrivateKey), nil
    }

    // otherwise its the path to a file containing the key
    return ioutil.ReadFile(p.PrivateKey)
}

// Returns the repository url for this project.
// A dummy username and password is inserted if no
// one is provided in the configuration file.
func (p *Project) GetRepository() (string) {
    repo := p.Repository

    // it's a http based repository url -> insert dummy
    // credentials if necessary
    if strings.HasPrefix(p.Repository, "http") ||
       strings.HasPrefix(p.Repository, "https") {
        repoUrl, err := url.Parse(p.Repository)
        if err != nil {
            return p.Repository
        }

        // we need to insert credentials
        if repoUrl.User == nil {
            repoUrl.User = url.UserPassword("anonymous", "null")
        }

        repo = repoUrl.String()
    }

    return repo
}
