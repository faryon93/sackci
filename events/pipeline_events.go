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

type EventFeed chan Event

type StageBegin struct {
    Stage string
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

type PipelineFinished struct {
    Status string
    Duration time.Duration
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (f EventFeed) StageBegin(stage string) {
    f <- StageBegin{stage}
}

func (f EventFeed) StageFinish(stage int, status string, duration time.Duration) {
    f <- StageFinish{stage, status, duration}
}

func (f EventFeed) StageLog(stage int, v ...interface{}) {
    message := strings.TrimSpace(fmt.Sprintln(v...))
    log.Info("pipeline", message)

    f <- StageLog{stage, message}
}

func (f EventFeed) PipelineFinished(status string, duration time.Duration) {
    f <- PipelineFinished{status, duration}
}
