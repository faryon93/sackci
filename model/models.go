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
    Id          uint64          `json:"id" storm:"id,increment"`
    Name        string          `json:"name"`
    BuildStatus string          `json:"status"`
    BuildId     int             `json:"build"`
    Time        time.Time       `json:"execution_time"`
    Duration    time.Duration   `json:"duration"`
}

type Env struct {
    Id      int    `json:"-" storm:"id,increment"`
    Project int    `json:"-" storm:"index"`
    Key     string `json:"key"`
    Value   string `json:"value"`
}

type Build struct {
    Id uint64 `json:"id"`
    Status string `json:"status"`
    Commit Commit `json:"commit"`
    Time time.Time `json:"time"`
    Duration string `json:"duration"`
    Node string `json:"node"`
    Stages []Stage `json:"stages"`
}

type Commit struct {
    Message string `json:"message"`
    Author string `json:"author"`
    Ref string `json:"ref"`
}

type Stage struct {
    Name string `json:"name"`
    Status string `json:"status"`
    Log []string `json:"log"`
}

