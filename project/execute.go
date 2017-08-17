package project
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
    "strings"
    "strconv"
    "errors"
    "path/filepath"

    "github.com/faryon93/sackci/agent"
    "github.com/faryon93/sackci/pipelinefile"
    "github.com/faryon93/sackci/model"
    "log"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    BUILDFILE = "Pipelinefile"
    WORKDIR = "/work"
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (p *Project) ExecuteBuild() {
    buildStart := time.Now()

    build := &model.Build{
        // TODO: propper id, num and node
        Project: 1,
        Status: model.BUILD_STATUS_RUNNING,
        Time: time.Now(),
        Num: 5,
        Node: "test",
        Commit: model.Commit{
            Message:"test",
            Author: "Maximilian Pachl",
            Ref: "904jf08qec9",
        },
    }
    scm := build.AddStage("SCM")
    build.Save()

    // construct the connection to the docker host
    agent := agent.Allocate()
    if agent == nil {
        addLog(build, scm, "no build agent is available, skipping build")
        return
    }
    defer agent.Free()

    // create the intermediate volume for the build pipeline
    start := time.Now()
    vol, err := agent.CreateVolume()
    if err != nil {
        addLog(build, scm, "failed to create volume:", err.Error())
        return
    }
    defer func() {
        agent.RemoveVolume(vol)
        addLog(build, scm, "removed temporary volume", hash(vol))
    }()
    addLog(build, scm, "temporary volume", hash(vol), "created in", time.Since(start))

    // checkout the latest revision to the build pipeline
    // TODO: support not only git scm
    start = time.Now()
    addLog(build, scm, "starting scm checkout for", p.Repository)
    container, ret, err := agent.Execute(vol, "sackci/git:latest", p.Repository)
    defer func() {
        if container != "" {
            agent.RemoveContainer(container)
            addLog(build, scm, "removed scm container", hash(container))
        }
    }()
    if err != nil {
        addLog(build, scm, "scm checkout failed:", err.Error())
        return
    }

    // make sure the scm checkout stage has completed successfully
    if ret != 0 {
        addLog(build, scm, "scm checkout exited with an error:", ret)
        return
    }
    addLog(build, scm, "scm checkout successfull finished in", time.Since(start))

    // download the Pipelinefile
    start = time.Now()
    pipline, err := getPipeline(agent, container)
    if err != nil {
        addLog(build, scm, "failed to download Pipelinefile:", err.Error())
        return
    }
    addLog(build, scm, "downloaded Pipelinefile from build volume in", time.Since(start))

    // print infos about the pipeline
    addLog(build, scm, "found", len(pipline.Stages), "stages", "(" + pipline.StageString() + ")", "in Pipelinefile")
    build.SetStageStatus(scm, model.STAGE_PASSED)

    // execute the configured stages
    for _, stage := range pipline.Stages {
        start := time.Now()
        stageId := build.AddStage(stage.Name)

        err := executeStage(agent, build, stageId, vol, &stage)
        if err != nil {
            addLog(build, scm, "failed to execute stage \"" + stage.Name + "\":", err.Error())
            build.SetStageStatus(stageId, model.STAGE_FAILED)
            return
        }

        addLog(build, stageId, "stage \"" + stage.Name + "\" finished successfully in", time.Since(start))
        build.SetStageStatus(stageId, model.STAGE_PASSED)
    }
    addLog(build, scm, "all pipeline stages finished")

    // save build artifacts if configured
    if len(pipline.Artifacts) > 0 {
        downloadArtifacts(agent, build, container, vol, pipline)
    }

    addLog(build, scm, "complete build finished in", time.Since(buildStart))

    build.Status = model.BUILD_STATUS_PASSED
    build.Save()
}


// ----------------------------------------------------------------------------------
//  private functions
// ----------------------------------------------------------------------------------

// Executes the given stage with the give volume.
func executeStage(agent *agent.Agent, build *model.Build, stageId int, vol string, stage *pipelinefile.Stage) (error) {
    // construct the build steps command string
    steps := strings.Join(stage.Steps, " && ")
    steps = "/bin/sh -c '" + steps + "'"

    // default to configured image if neceasarry
    image := "alpine:latest"
    if len(stage.Image) > 0 {
        image = stage.Image
    }

    // execute the build steps in a new container
    addLog(build, stageId, "executing stage \"" + stage.Name + "\" with image \"" + image + "\"")
    id, ret, err := agent.Execute(vol, image, steps)
    defer func() {
        if id != "" {
            agent.RemoveContainer(id)
            addLog(build, stageId, "removed stage container", hash(id))
        }
    }()
    if err != nil {
        return err
    }

    // check the return code
    if ret != 0 {
        return errors.New("error code " + strconv.Itoa(ret))
    }

    return nil
}

// Return the projects parsed Pipelinefile.
func getPipeline(agent *agent.Agent, container string) (*pipelinefile.Pipelinefile, error) {
    // download the pipeline file
    pipelineContent, err := agent.GetFile(container, BUILDFILE)
    if err != nil {
        return nil, err
    }

    // parse the pipelinefile
    pipline, err := pipelinefile.Parse(pipelineContent)
    if err != nil {
        return nil, err
    }

    return pipline, nil
}

func downloadArtifacts(agent *agent.Agent, build *model.Build, container string, vol string, pipline *pipelinefile.Pipelinefile) {
    start := time.Now()

    // save as tar.gz file to the local filesystem
    file := filepath.Join("E:/tmp/", vol + ".tar.gz")
    err := agent.SaveFile(container, agent.Filepath(WORKDIR, pipline.Artifacts), file)
    if err != nil {
        addLog(build, 0, "failed to download build artifacts:", err.Error())
        return
    }

    addLog(build, 0, "successfully downloaded artifact archive in", time.Since(start))
}

func addLog(b *model.Build, stage int, v ...interface{}) {
    log.Println(v)
    b.Log(stage, v)
}

// Gets the short representation of an full length hash.
func hash(container string) (string) {
    return container[0:12]
}
