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
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type EventBase struct {
    Project int `json:"project_id"`
    Build int `json:"build_num"`
    Timestamp int64 `json:"timestamp"`
}

type EvtStageBegin struct {
    *EventBase
    Stage int `json:"stage"`
    Status string `json:"status"`
}

type EvtStageFinish struct {
    *EventBase
    Stage int `json:"stage"`
    Status string `json:"status"`
    Duration time.Duration `json:"duration"`
}

type EvtStageLog struct {
    *EventBase
    Stage int `json:"stage"`
    Message string `json:"message"`
}

type EvtPipelineBegin struct {
    *EventBase
    Time time.Time `json:"time"`
    Status string `json:"status"`
}

type EvtPipelineFinished struct {
    *EventBase
    Status string `json:"status"`
    Duration time.Duration `json:"duration"`
}

type EvtPipelineFound struct {
    *EventBase
    Stages []string `json:"stages"`
}

type EvtCommitFound struct {
    *EventBase
    Commit Commit `json:"commit"`
}
