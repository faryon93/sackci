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
    "github.com/faryon93/sackci/events"
    "github.com/faryon93/sackci/pipelinefile"
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
    Events events.EventFeed

    // private variables
    mutex sync.Mutex
    project *model.Project
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
        Events: make(events.EventFeed, EVENT_STREAM_BUFFER),
    }, nil
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (p *Pipeline) SetProject(project *model.Project) (error) {
    // assign the project to the pipeline
    if p.project != nil {
        return ErrAlreadyExecuted
    }

    p.project = project

    return nil
}

