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
    "time"
    "fmt"
    "strings"

    "github.com/faryon93/sackci/log"
)

// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

type EvtStageBegin struct {
    Stage int
}

type EvtStageFinish struct {
    Stage int
    Status string
    Duration time.Duration
}

type EvtStageLog struct {
    Stage int
    Message string
}

type EvtPipelineBegin struct {
    Time time.Time
}

type EvtPipelineFinished struct {
    Status string
    Duration time.Duration
}

type EvtPipelineFound struct {
    Stages []string
}

type EvtCommitFound struct {
    Commit Commit
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (f EventFeed) StageBegin(stage int) {
    f <- EvtStageBegin{stage}
}

func (f EventFeed) StageFinish(stage int, status string, duration time.Duration) {
    f <- EvtStageFinish{stage, status, duration}
}

func (f EventFeed) StageLog(stage int, v ...interface{}) {
    message := "\u001b[0;33m[pipeline] " + strings.TrimSpace(fmt.Sprintln(v...)) + "\u001b[m"
    log.Info("pipeline", message)

    f <- EvtStageLog{stage, message}
}

func (f EventFeed) ConsoleLog(stage int, v ...interface{}) {
    message := fmt.Sprintln(v...)
    log.Info("pipeline", strings.TrimSpace(message))

    f <- EvtStageLog{stage, message}
}

func (f EventFeed) PipelineFinished(status string, duration time.Duration) {
    f <- EvtPipelineFinished{status, duration}
}

func (f EventFeed) PipelineFound(stages []string) {
    f <- EvtPipelineFound{stages}
}

func (f EventFeed) PipelineBegin(time time.Time) {
    f <- EvtPipelineBegin{time}
}

func (f EventFeed) CommitFound(commit *Commit) {
    f <- EvtCommitFound{Commit: *commit}
}
