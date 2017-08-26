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
    // common constants
    LOG_TAG = "scm"

    // minimal polling interval in seconds
    MIN_POLLING_INTERVAL = 10
)


// --------------------------------------------------------------------------------------
//  global variables
// --------------------------------------------------------------------------------------

var (
    timers = []*util.CycleTimer{}
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
func poll(project *model.Project, timer *util.CycleTimer) {
    // etry and exit logs
    log.Info(LOG_TAG, "setup project \"" + project.Name + "\" for scm polling (interval:", project.Interval, ")")
    defer func() {
        log.Info(LOG_TAG, "scm polling for project \"" + project.Name + "\" exited")
        waitgroup.Done()
    }()

    // some runtime variables
    oldRef := ""

    // cycle
    for {
        // if an error on the timer occours -> finish
        err := timer.Wait(time.Duration(project.Interval) * time.Second)
        if err == util.ErrTimerCancel {
            return

        } else if err != nil {
            log.Error(LOG_TAG, "scm cycle timer for project failed:", err.Error())
            return
        }

        // begin of the cycle
        start := time.Now()
        log.Info(LOG_TAG, "starting scm polling for project \"" + project.Name + "\"")

        // create a new pipeline
        pipeline, err := agent.CreatePipeline()
        if err != nil {
            log.Error(LOG_TAG, "failed to create scm polling pipeline:", err.Error())
            continue
        }
        pipeline.SetProject(project)

        // check if changes have happend since the last polling cycle
        newRef, err := pipeline.HeadRef()
        if err != nil {
            log.Error(LOG_TAG, "failed to compare scm refs:", err.Error())
            continue
        }

        // a new reference is available -> execute the build
        if oldRef != newRef {
            log.Info(LOG_TAG, "a new ref", newRef, "for project \"" + project.Name + "\" is available")
            oldRef = newRef
        } else {
            log.Info(LOG_TAG, "project", project.Name, "is up to date")
        }

        // TODO: execute the pipeline

        // cleanup the pipeline
        pipeline.Destroy()

        log.Info(LOG_TAG, "scm polling took", time.Since(start))
    }
}