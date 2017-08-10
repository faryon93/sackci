package model

import (
    "time"
    "github.com/boltdb/bolt"
    "encoding/json"
)

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

    STAGE_FAILED = "failed"
    STAGE_PASSED = "passed"
    STAGE_IGNORED = "ignored"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

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


// --------------------------------------------------------------------------------------
//  public members
// --------------------------------------------------------------------------------------

func (b *Build) SetId(id uint64) {
    b.Id = id
}


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func GetProjectHistory(project uint64) ([]Build, error) {
    builds := make([]Build, 0)

    // find all build indicies
    var buildIndex []uint64
    err := Get(BUILDS_INDEX_BUCKET, project, &buildIndex)
    if err != nil {
        return builds, err
    }

    return builds, db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(BUILDS_BUCKET))

        for _, buildId := range buildIndex {
            buf := b.Get(itob(buildId))
            if buf == nil {
                continue
            }

            var build Build
            err := json.Unmarshal(buf, &build)
            if err != nil {
                return err
            }

            builds = append(builds, build)
        }

        return nil
    })

    return builds, nil
}
