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
    Stage int `json:"stage"`
    Status string `json:"status"`
}

func (e *EvtStageBegin) Event() string {
    return "stage_begin"
}

type EvtStageFinish struct {
    Stage int `json:"stage"`
    Status string `json:"status"`
    Duration time.Duration `json:"duration"`
}

func (e *EvtStageFinish) Event() string {
    return "stage_finish"
}

type EvtStageLog struct {
    Stage int `json:"stage"`
    Message string `json:"message"`
}

func (e *EvtStageLog) Event() string {
    return "stage_log"
}

type EvtPipelineBegin struct {
    Time time.Time `json:"time"`
    Status string `json:"status"`
}

func (e *EvtPipelineBegin) Event() string {
    return "pipeline_begin"
}

type EvtPipelineFinished struct {
    Status string `json:"status"`
    Duration time.Duration `json:"duration"`
}

func (e *EvtPipelineFinished) Event() string {
    return "pipeline_finish"
}

type EvtPipelineFound struct {
    Stages []string `json:"stages"`
}

func (e *EvtPipelineFound) Event() string {
    return "pipeline_found"
}

type EvtCommitFound struct {
    Commit Commit `json:"commit"`
}

func (e *EvtCommitFound) Event() string {
    return "commit_found"
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (f EventFeed) StageBegin(stage int) {
    f <- &EvtStageBegin{stage, STAGE_RUNNING}
}

func (f EventFeed) StageFinish(stage int, status string, duration time.Duration) {
    f <- &EvtStageFinish{stage, status, duration}
}

func (f EventFeed) StageLog(stage int, v ...interface{}) {
    trimmed := strings.TrimSpace(fmt.Sprintln(v...))

    log.Info("pipeline", trimmed)
    message := "\u001b[0;33m[pipeline] " + trimmed + "\u001b[m"
    f <- &EvtStageLog{stage, message}
}

func (f EventFeed) ConsoleLog(stage int, v ...interface{}) {
    f <- &EvtStageLog{stage, fmt.Sprintln(v...)}
}

func (f EventFeed) PipelineFinished(status string, duration time.Duration) {
    f <- &EvtPipelineFinished{status, duration}
}

func (f EventFeed) PipelineFound(stages []string) {
    f <- &EvtPipelineFound{stages}
}

func (f EventFeed) PipelineBegin(time time.Time) {
    f <- &EvtPipelineBegin{time, BUILD_RUNNING}
}

func (f EventFeed) CommitFound(commit *Commit) {
    f <- &EvtCommitFound{Commit: *commit}
}
