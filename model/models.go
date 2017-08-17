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
//  constants
// --------------------------------------------------------------------------------------

const (
    BUILDS_BUCKET = "builds"
    BUILDS_INDEX_BUCKET = "builds_index"

    BUILD_STATUS_PASSED = "passed"
    BUILD_STATUS_FAILED = "failed"
    BUILD_STATUS_RUNNING = "running"
    BUILD_STATUS_WAITING = "waiting"

    STAGE_FAILED = "failed"
    STAGE_PASSED = "passed"
    STAGE_IGNORED = "ignored"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type Env struct {
    Id      int    `json:"-" storm:"id,increment"`
    Project int    `json:"-" storm:"index"`
    Key     string `json:"key" groups:"queryall"`
    Value   string `json:"value" groups:"queryall"`
}
