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
    "time"
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Publishes all information about the Pipeline definition
// to the event stream.
func (p *Pipeline) PublishPipeline() {
    if p.definition == nil {
        return
    }

    // get all stage names form the Pipeline definition
    stages := make([]string, len(p.definition.Stages))
    for i, stage := range p.definition.Stages {
        stages[i] = stage.Name
    }

    // publish the event
    p.Events.PipelineFound(stages)
}

// Appends a log file to a stage.
func (p *Pipeline) Log(stage int, v ...interface{}) {
    p.Events.StageLog(stage, v...)
}

// Begins the given stage.
func (p *Pipeline) BeginStage(stage int) {
    p.Events.StageBegin(stage)
}

// Finishes the given stage
func (p *Pipeline) FinishStage(stage int, status string, duration time.Duration) {
    p.Events.StageFinish(stage, status, duration)
}

// Inform about begining a pipeline.
func (p *Pipeline) BeginPipeline(start time.Time) {
    p.Events.PipelineBegin(start)
}

// Inform about the end of a pipeline.
func (p *Pipeline) FinishPipeline(status string, duration time.Duration) {
    p.Events.PipelineFinished(status, duration)
}