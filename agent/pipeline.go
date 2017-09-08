package agent
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
    "time"

    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/pipelinefile"
    "github.com/faryon93/sackci/log"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    EVENT_STREAM_BUFFER = 128
    LOG_TAG = "pipeline"
)

// errors
var (
    ErrNoContainer = errors.New("no container available")
    ErrNoAgent = errors.New("no agent available")
    ErrAlreadyExecuted = errors.New("pipeline was already executed")
    ErrInvalidBranch = errors.New("branch does not exist")
    ErrNoProject = errors.New("no project set")
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type Pipeline struct {
    // public variables
    Agent *Agent
    Volume string
    Containers []string
    StartTime time.Time
    Env map[string]string
    Events chan interface{}

    // private variables
    mutex sync.Mutex
    project *model.Project
    build *model.Build
    definition *pipelinefile.Definition
}


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

// Creates a new pipeline on a free agent.
func CreatePipeline() (*Pipeline, error) {
    start := time.Now()

    agent := Allocate()
    if agent == nil {
        return nil, ErrNoAgent
    }

    volume, err := agent.CreateVolume()
    if err != nil {
        return nil, err
    }

    return &Pipeline{
        Agent: agent,
        Volume: volume,
        Containers: []string{},
        StartTime: start,
        Env: make(map[string]string),
        Events: make(chan interface{}, EVENT_STREAM_BUFFER),
    }, nil
}

// Destroys the whole pipeline
func (p *Pipeline) Destroy() {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    // remove all containers
    for _, container := range p.Containers {
        err := p.Agent.RemoveContainer(container)
        if err != nil {
            log.Error(LOG_TAG, "failed to remove container:", err.Error())
            continue
        }
    }

    // destroy the volume
    err := p.Agent.RemoveVolume(p.Volume)
    if err != nil {
        log.Error(LOG_TAG, "failed to remove volume:", err.Error())
    }

    // free the agent
    p.Agent.Free()

    // close the event stream
    close(p.Events)
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Assigns a project to this pipeline.
func (p *Pipeline) SetProject(project *model.Project) {
    p.project = project

    // copy the project env
    if project.Env != nil {
        for key, val := range project.Env {
            p.Env[key] = val
        }
    }
}

// Assigns a build to this pipeline.
func (p *Pipeline) SetBuild(build *model.Build) {
    p.build = build
}


// ----------------------------------------------------------------------------------
//  private members
// ----------------------------------------------------------------------------------

// Gets the accumulated environment (pipeline and project) settings of this pipeline.
func (p *Pipeline) getEnv() []string {
    env := make([]string, len(p.Env) + len(p.project.Env))
    i := 0

    // project specific variables
    for key, val := range p.project.Env {
        env[i] = key + "=" + val
        i++
    }

    // pipeline specific variables
    for key, val := range p.Env {
        env[i] = key + "=" + val
        i++
    }

    return env
}

// Adds a container for destruction to this pipeline.
func (p *Pipeline) addContainer(container string) {
    p.mutex.Lock()
    p.Containers = append(p.Containers, container)
    p.mutex.Unlock()
}
