package agent
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
    "errors"
    "strconv"
    "time"

    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/pipelinefile"
    "github.com/faryon93/sackci/log"
    "github.com/faryon93/sackci/util"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    PIPELINEFILE = "Pipelinefile"
    STAGE_IMAGE = "alpine:latest"
    KEY_PATH = "/tmp/id_rsa"

    // scm stage definitions
    STAGE_SCM_ID = 0
    STAGE_SCM_NAME = "Clone"

    // internal environment variables
    CI_BUILDNR = "CI_BUILDNR"
    CI_AGENT = "CI_AGENT"
    CI_COMMIT_REF = "CI_COMMIT_REF"
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Executes the give project on this pipeline.
// This is a oneshot action. The pipeline is toren down on success or error.
func (p *Pipeline) Execute() (error) {
    if p.project == nil {
        return ErrNoProject
    }

    // begin the build for the project
    log.Info(LOG_TAG,"executing build for project \"" + p.project.Name + "\"")
    p.BeginPipeline(p.StartTime, p.Agent.Name)

    // assign the ci server generated environment variables
    p.Env[CI_BUILDNR] = strconv.Itoa(p.build.Num)
    p.Env[CI_AGENT] = p.Agent.Name

    // get a working copy of the repo
    scmStart := time.Now()
    p.BeginStage(STAGE_SCM_ID)
    p.Log(STAGE_SCM_ID,"starting scm checkout for", util.MaskCredentials(p.project.Repository))
    commit, err := p.Clone()
    if err != nil {
        p.Log(STAGE_SCM_ID,"scm checkout failed:", err.Error())
        p.FinishStage(STAGE_SCM_ID, model.STAGE_FAILED, time.Since(scmStart))
        p.FinishPipeline(model.BUILD_FAILED, time.Since(p.StartTime))
        return err
    }
    p.CommitFound(commit)
    p.Env[CI_COMMIT_REF] = commit.Ref
    p.Log(STAGE_SCM_ID, "scm checkout completed successfully in", time.Since(scmStart))

    // get the pipeline definition
    start := time.Now()
    definition, err := p.GetPipelinefile()
    if err != nil {
        p.Log(STAGE_SCM_ID, "failed to get Pipelinefile:", err.Error())
        p.FinishStage(STAGE_SCM_ID, model.STAGE_FAILED, time.Since(start))
        p.FinishPipeline(model.BUILD_FAILED, time.Since(p.StartTime))
        return err
    }
    p.definition = definition

    // the SCM stage has sucessfully finished
    p.Log(STAGE_SCM_ID, "sucessfully obtained Pipelinefile in", time.Since(start))
    p.Log(STAGE_SCM_ID, "found", len(definition.Stages), "stages", "(" + definition.StageString() + ") in Pipelinefile")
    p.PublishPipeline()
    p.FinishStage(STAGE_SCM_ID, model.STAGE_PASSED, time.Since(scmStart))

    // execute all configured stages
    for stageId, stage := range definition.Stages {
        start := time.Now()

        // the "Prolog" stage is added automatically
        stageId = stageId + 1

        // execute the stage
        err := p.ExecuteStage(stageId, &stage)
        if err != nil {
            p.Log(stageId, "stage \"" + stage.Name + "\" failed:", err.Error(), "in", time.Since(start))
            p.FinishStage(stageId, model.STAGE_FAILED, time.Since(start))
            p.FinishPipeline(model.BUILD_FAILED, time.Since(p.StartTime))
            return err
        }

        // stage executed successfully
        p.Log(stageId, "stage \"" + stage.Name + "\" completed in", time.Since(start))
        p.FinishStage(stageId, model.STAGE_PASSED, time.Since(start))
    }

    p.FinishPipeline(model.BUILD_PASSED, time.Since(p.StartTime))

    return nil
}

// Downloads and parses the pipeline defition from the working copy.
func (p *Pipeline) GetPipelinefile() (*pipelinefile.Definition, error) {
    // download content of the Pipelinefile
    content, err := p.ReadFile(PIPELINEFILE)
    if err != nil {
        return nil, err
    }

    // parse the Pipelinefile
    return pipelinefile.Parse(content)
}

// Executes the given stage on this Pipieline.
func (p *Pipeline) ExecuteStage(stageId int, stage *pipelinefile.Stage) (error) {
    // construct the build steps command string
    // insert an echo of the command
    steps := "/bin/sh -c '"
    for i := 0; i < len(stage.Steps); i++ {
        steps += "echo $: " + stage.Steps[i] + " && " + stage.Steps[i]
        if i < len(stage.Steps) - 1 {
            steps += " && "
        }
    }
    steps += "'"

    // default to configured image if neceasarry
    image := STAGE_IMAGE
    if len(stage.Image) > 0 {
        image = stage.Image
    }

    // begin the stage
    p.BeginStage(stageId)
    p.Log(stageId, "executing stage \"" + stage.Name + "\"", "in image \"" + image + "\"")

    // execute the steps inside a container on the build agent
    ret, err := p.Container(image, steps, WORKDIR, func(line string) {
        p.LogTerminal(stageId, util.MaskCredentials(line))
    })
    if err != nil {
        return err
    }

    if ret != 0 {
        return errors.New("error code " + strconv.Itoa(ret))
    }

    return nil
}
