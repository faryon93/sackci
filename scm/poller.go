package scm
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

    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/agent"
    "github.com/faryon93/sackci/log"
    "github.com/faryon93/sackci/ctx"
    "github.com/faryon93/sackci/util"
    "sync"
)


// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    // minimal polling interval in seconds
    MIN_POLLING_INTERVAL = 10
)


// --------------------------------------------------------------------------------------
//  global variables
// --------------------------------------------------------------------------------------

var (
    timers = []*util.Timer{}
    waitgroup = sync.WaitGroup{}
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

// Setup the necesarry go routines for all projects with "polling" trigger.
func Setup() {
    for i := 0; i < len(ctx.Conf.Projects); i++ {
        project := &ctx.Conf.Projects[i]

        if project.Trigger == model.TRIGGER_POLL {
            // ensure the minimal polling interval
            if project.Interval < MIN_POLLING_INTERVAL {
                project.Interval = MIN_POLLING_INTERVAL
            }

            // construct the timer to be used
            timer := util.NewTimer()

            // add to the scm pool
            waitgroup.Add(1)
            timers = append(timers, timer)

            // execute the poller async
            go poll(project, timer)

        }
    }
}

func Destroy() {
    // cancel all timers
    for _, timer := range timers {
        timer.Cancel()
    }

    // wait for all polling routines to stop
    waitgroup.Wait()
}


// --------------------------------------------------------------------------------------
//  private functions
// --------------------------------------------------------------------------------------

// The scm polling loop.
func poll(project *model.Project, timer *util.Timer) {
    // etry and exit logs
    log.Info("scm", "setup project \"" + project.Name + "\" for scm polling (interval:", project.Interval, ")")
    defer func() {
        log.Info("scm", "scm polling for project \"" + project.Name + "\" exited")
        waitgroup.Done()
    }()

    // some runtime variables
    oldRef := ""

    // cycle
    for {
        // if an error on the timer occours -> finish
        err := timer.Wait(time.Duration(project.Interval) * time.Second)
        if err == util.ErrCancel {
            return

        } else if err != nil {
            log.Error("scm", "scm cycle timer for project failed:", err.Error())
            return
        }

        // begin of the cycle
        start := time.Now()
        log.Info("scm", "starting scm polling for project \"" + project.Name + "\"")

        // create a new pipeline
        pipeline, err := agent.CreatePipeline()
        if err != nil {
            log.Error("scm", "failed to create scm polling pipeline:", err.Error())
            continue
        }
        pipeline.SetProject(project)

        // check if changes have happend since the last polling cycle
        repoChanged, ref, err := pipeline.Compare(oldRef)
        if err != nil {
            log.Error("scm", "failed to compare scm refs:", err.Error())
            continue
        }

        // a new reference is available -> execute the build
        if repoChanged {
            log.Info("scm", "a new ref", ref, "for project \"" + project.Name + "\" is available")
            oldRef = ref

        } else {
            log.Info("scm", "project \"" + project.Name + "\" is up to date")
        }

        // TODO: execute the pipeline

        // cleanup the pipeline
        pipeline.Destroy()

        log.Info("scm", "scm polling took", time.Since(start))
    }
}