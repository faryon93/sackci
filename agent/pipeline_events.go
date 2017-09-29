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
    "strings"
    "fmt"

    log "github.com/sirupsen/logrus"

    "github.com/faryon93/sackci/model"
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
    p.Events <- model.EvtPipelineFound{
        p.getBaseEvent(),stages,
    }
}

// Publishes the commit details.
func (p *Pipeline) CommitFound(commit *model.Commit) {
    p.Events <- model.EvtCommitFound{
        p.getBaseEvent(), *commit,
    }
}

// Appends a log line to a stage.
func (p *Pipeline) Log(stage int, v ...interface{}) {
    trimmed := strings.TrimSpace(fmt.Sprintln(v...))
    message := "\u001b[0;33m[pipeline] " + trimmed + "\u001b[m\n"

    log.Infoln(trimmed)
    p.Events <- model.EvtStageLog{
        p.getBaseEvent(),
        stage, message,
    }
}

// Appends a terminal log line to a stage.
func (p *Pipeline) LogTerminal(stage int, v ...interface{}) {
    msg := fmt.Sprintln(v...)
    if msg == "\n" {
        return
    }

    p.Events <- model.EvtStageLog{
        p.getBaseEvent(),
        stage, msg,
    }
}

// Begins the given stage.
func (p *Pipeline) BeginStage(stage int) {
    p.Events <- model.EvtStageBegin{
        p.getBaseEvent(),
        stage, model.STAGE_RUNNING,
    }
}

// Finishes the given stage
func (p *Pipeline) FinishStage(stage int, status string, duration time.Duration) {
    p.Events <- model.EvtStageFinish{
        p.getBaseEvent(),
        stage, status, duration,
    }
}

// Inform about begining a pipeline.
func (p *Pipeline) BeginPipeline(start time.Time, agent string) {
    p.Events <- model.EvtPipelineBegin{
        p.getBaseEvent(),
        start, model.BUILD_RUNNING, agent,
    }
}

// Inform about the end of a pipeline.
func (p *Pipeline) FinishPipeline(status string, duration time.Duration) {
    p.Events <- model.EvtPipelineFinished{
        p.getBaseEvent(),
        status, duration,
    }
}


// ----------------------------------------------------------------------------------
//  private members
// ----------------------------------------------------------------------------------

func (p *Pipeline) getBaseEvent() (*model.EventBase) {
    project := 0
    if p.project != nil {
        project = p.project.Id
    }

    build := 0
    if p.build != nil {
        build = int(p.build.Num)
    }

    return &model.EventBase{
        Project: project,
        Build: build,
        Timestamp: time.Now().UnixNano(),
    }
}
