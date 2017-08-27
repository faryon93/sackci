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

    // if a project is assigned
    // we need to clear the execution lock
    if p.project != nil {
        p.project.ExecutionLock.Unlock()
    }
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Assigns a project to this pipeline.
// After this time the execution lock of the project is locked.
// As recently as the pipeline is destroyed the lock is returned.
func (p *Pipeline) SetProject(project *model.Project) {
    p.project = project
    p.project.ExecutionLock.Lock()
}

func (p *Pipeline) SetBuild(build *model.Build) {
    p.build = build
}
