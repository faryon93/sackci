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
    "sync"

    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/agent"
    "github.com/faryon93/sackci/log"
    "github.com/faryon93/sackci/ctx"
    "github.com/faryon93/sackci/util"
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

// Setup the necessary go routines for all projects with "polling" trigger.
func Setup() {
    for i := 0; i < len(ctx.Conf.Projects); i++ {
        project := &ctx.Conf.Projects[i]

        // we are not interested in non-polling projects
        if project.Trigger != model.TRIGGER_POLL {
            continue
        }

        // ensure the minimal polling interval
        if project.Interval < MIN_POLLING_INTERVAL {
            project.Interval = MIN_POLLING_INTERVAL
        }

        // construct the timer to be used
        interval := time.Duration(project.Interval) * time.Second
        timer := util.NewTimer(interval, func(t *util.CycleTimer) {
            poll(project, t)
        }, func() {
            log.Info(LOG_TAG, "polling for \"" + project.Name + "\" exited")
            waitgroup.Done()
        })

        // add to the scm pool
        waitgroup.Add(1)
        timers = append(timers, timer)

        // execute the poller async
        log.Info(LOG_TAG, "setup project \"" + project.Name +
                             "\" for scm polling (interval:", project.Interval, ")")
        timer.Start()
    }
}

// Destroys all scm polling goroutines.
// Blocks until all routines have exited.
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
func poll(project *model.Project, t *util.CycleTimer) {
    // create a new pipeline
    pipeline, err := agent.CreatePipeline()
    if err != nil {
        log.Error(LOG_TAG, "failed to create scm polling pipeline:", err.Error())
        return
    }
    pipeline.SetProject(project)
    pipeline.SetArtifactsDir(ctx.Conf.GetArtifactsDir())
    defer pipeline.Destroy()    // make sure the pipeline gets destroyed

    log.Info(LOG_TAG, "starting scm polling for project \"" + project.Name + "\"")

    // check if changes have happend since the last polling cycle
    newRef, err := pipeline.HeadRef()
    if err != nil {
        log.Error(LOG_TAG, "failed to compare scm refs:", err.Error())
        return
    }

    // get the last build for the project
    lastBuild, err := project.GetLastBuild()
    if err != nil {
        log.Error(LOG_TAG, "failed to get get last build:", err.Error())
        return
    }

    log.Info(LOG_TAG, "scm polling took", time.Since(t.StartTime))

    // detect if new changes are available in the repository
    // trigger a new build if changes are detected
    if lastBuild == nil || lastBuild.Commit.Ref != newRef {
        // try to lock the project, if a build is already running
        // skip this build cycle
        if ok := project.Lock(); ok != nil {
            return
        }
        defer project.Unlock()

        log.Error(LOG_TAG, "changes detected for project", project.Name,
                              "with ref:", util.ShortHash(newRef))

        // construct the build object for saving in the database
        build := project.NewBuild()
        build.Time = time.Now()
        build.Node = pipeline.Agent.Name
        build.Stages = []model.Stage{
            {Name: agent.STAGE_SCM_NAME, Status: model.STAGE_IGNORED, Log: []string{}},
        }
        err = build.Save()
        if err != nil {
            log.Error(LOG_TAG, "failed to save build:", err.Error())
            return
        }
        pipeline.SetBuild(build)

        // redirect all events for the eventstream
        go func() {
            for event := range pipeline.Events {
                build.Publish(event)
                ctx.Feed.Publish(event)
            }
        }()

        // execute the pipeline
        err = pipeline.Execute()
        if err != nil {
            log.Error(LOG_TAG, "failed to execute pipeline:", err.Error())
            return
        }
    }
}
