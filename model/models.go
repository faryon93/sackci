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

// --------------------------------------------------------------------------------------
//  imports
// --------------------------------------------------------------------------------------

import (
    "time"
)


// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    BUILDS_BUCKET = "builds"
    BUILDS_INDEX_BUCKET = "builds_index"

    BUILD_STATUS_PASSED = "passed"
    BUILD_STATUS_FAILED = "failed"
    BUILD_STATUS_RUNNING = "running"

    STAGE_FAILED = "failed"
    STAGE_PASSED = "passed"
    STAGE_IGNORED = "ignored"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type Project struct {
    Id          uint64          `json:"id" storm:"id,increment" groups:"all,one"`
    Name        string          `json:"name" groups:"all,one"`
    BuildStatus string          `json:"status" groups:"all,one"`
    BuildId     int             `json:"build" groups:"all,one"`
    Time        time.Time       `json:"execution_time" groups:"all,one"`
    Duration    time.Duration   `json:"duration" groups:"all,one"`
}

type Env struct {
    Id      int    `json:"-" storm:"id,increment"`
    Project int    `json:"-" storm:"index"`
    Key     string `json:"key" groups:"queryall"`
    Value   string `json:"value" groups:"queryall"`
}

type Build struct {
    Id          uint64      `json:"id" storm:"id,increment" groups:"queryall,one"`
    Project     uint64      `json:"-" strom:"index"`

    Num         uint64          `json:"num" strom:"index" groups:"queryall,one"`
    Status      string          `json:"status" groups:"queryall,one"`
    Commit      Commit          `json:"commit" groups:"queryall,one"`
    Time        time.Time       `json:"time" groups:"queryall,one"`
    Duration    time.Duration   `json:"duration" groups:"queryall,one"`
    Node        string          `json:"node" groups:"queryall,one"`
    Stages      []Stage         `json:"stages" groups:"one"`
}

type Commit struct {
    Message string `json:"message" groups:"queryall,one"`
    Author  string `json:"author" groups:"queryall,one"`
    Ref     string `json:"ref" groups:"queryall,one"`
}

type Stage struct {
    Name    string      `json:"name" groups:"one"`
    Status  string      `json:"status" groups:"one"`
    Log     []string    `json:"log" groups:"one"`
}
