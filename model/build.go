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
    "github.com/faryon93/sackci/events"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type Build struct {
    Id          uint64          `json:"id" storm:"id,increment" groups:"queryall,one"`
    Project     uint64          `json:"-" strom:"index"`
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


// --------------------------------------------------------------------------------------
//  public members
// --------------------------------------------------------------------------------------

func (b *Build) AddStage(stage string) (int) {
    b.Stages = append(b.Stages, Stage{
        Name: stage,
        Status: STAGE_IGNORED,
        Log: []string{},
    })

    return len(b.Stages) - 1
}

func (b *Build) Log(stage int, message string) {
    if stage >= len(b.Stages) {
        return
    }

    b.Stages[stage].Log = append(b.Stages[stage].Log, message)
    b.Save()
}

func (b *Build) Save() (error) {
    if b.Id == 0 {
        return Get().Save(b)
    } else {
        return Get().Update(b)
    }
}

func (b *Build) Attach(src chan events.Event) {
    for event := range src {
        if evt, ok := event.(events.StageBegin); ok {
            b.Stages = append(b.Stages, Stage{
                Name: evt.Stage,
                Status: STAGE_IGNORED,
                Log: []string{},
            })

        } else if evt, ok := event.(events.StageFinish); ok {
            if evt.Stage >= len(b.Stages) {
                continue
            }

            b.Stages[evt.Stage].Status = evt.Status
            // TODO: duration

        } else if evt, ok := event.(events.StageLog); ok {
            if evt.Stage >= len(b.Stages) {
                continue
            }

            b.Stages[evt.Stage].Log = append(b.Stages[evt.Stage].Log, evt.Message)

        } else if evt, ok := event.(events.PipelineFinished); ok {
            b.Status = evt.Status
            b.Duration = evt.Duration
        }

        b.Save()
    }
}