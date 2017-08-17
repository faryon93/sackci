package events
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

type StageBegin struct {
    Stage int
}

type StageFinish struct {
    Stage int
    Status string
    Duration time.Duration
}

type StageLog struct {
    Stage int
    Message string
}

type PipelineBegin struct {
    Time time.Time
}

type PipelineFinished struct {
    Status string
    Duration time.Duration
}

type PipelineFound struct {
    Stages []string
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (f EventFeed) StageBegin(stage int) {
    f <- StageBegin{stage}
}

func (f EventFeed) StageFinish(stage int, status string, duration time.Duration) {
    f <- StageFinish{stage, status, duration}
}

func (f EventFeed) StageLog(stage int, v ...interface{}) {
    message := "\u001b[0;33m[pipeline] " + strings.TrimSpace(fmt.Sprintln(v...)) + "\u001b[m"
    log.Info("pipeline", message)

    f <- StageLog{stage, message}
}

func (f EventFeed) ConsoleLog(stage int, v ...interface{}) {
    message := fmt.Sprintln(v...)
    log.Info("pipeline", strings.TrimSpace(message))

    f <- StageLog{stage, message}
}

func (f EventFeed) PipelineFinished(status string, duration time.Duration) {
    f <- PipelineFinished{status, duration}
}

func (f EventFeed) PipelineFound(stages []string) {
    f <- PipelineFound{stages}
}

func (f EventFeed) PipelineBegin(time time.Time) {
    f <- PipelineBegin{time}
}