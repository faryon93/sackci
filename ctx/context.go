package ctx
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
//  types
// --------------------------------------------------------------------------------------

import (
    "time"

    "github.com/faryon93/sackci/sse"
)

// --------------------------------------------------------------------------------------
//  global variables
// --------------------------------------------------------------------------------------

var (
    Feed *sse.Group
)


// --------------------------------------------------------------------------------------
//  initializer
// --------------------------------------------------------------------------------------

func Init() {
    // initialize the context variables
    Feed = sse.NewFeed()

    // event stream for testing purpose
    go func() {
        status := "failed"

        for {
            if status == "failed" {
                status = "running"
            } else {
                status = "failed"
            }
            Feed.Publish(&test{1, status})

            time.Sleep(2 * time.Second)
        }
    }()
}


// --------------------------------------------------------------------------------------
//  intermediate types
// --------------------------------------------------------------------------------------

type test struct {
    Project int `json:"project"`
    Status string `json:"status"`
}

func (t *test) Event() string {
    return "project_changed"
}
