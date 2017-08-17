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
    "strings"

    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/pipelinefile"
    "github.com/faryon93/sackci/log"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    SCM_IMAGE = "sackci/git:latest"
    STAGE_IMAGE = "alpine:latest"
    PIPELINEFILE = "Pipelinefile"
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Executes the give project on this pipeline.
// This is a oneshot action. The pipeline is toren down on success or error.
func (p *Pipeline) Execute(project *model.Project) (error) {
    // assign the project to the pipeline
    if p.project != nil {
        return ErrAlreadyExecuted
    }
    p.project = project

    // whenever we exit this funtion -> destroy the whole pipeline
    defer p.Destroy()

    log.Info(LOG_TAG,"executing build for project \"" + project.Name + "\"")

    // get a working copy of the repo
    start := time.Now()
    log.Info(LOG_TAG,"starting scm checkout for", project.Repository)
    err := p.Checkout()
    if err != nil {
        log.Error(LOG_TAG,"scm checkout failed:", err.Error())
        return err
    }
    log.Info(LOG_TAG,"scm checkout completed successfully in", time.Since(start))

    // get the pipeline definition
    start = time.Now()
    definition, err := p.GetPipelinefile()
    if err != nil {
        log.Info(LOG_TAG,"failed to get Pipelinefile:", err.Error())
    }
    log.Info(LOG_TAG,"sucessfully obtained Pipelinefile in", time.Since(start))
    log.Info(LOG_TAG,"found", len(definition.Stages), "stages", "(" + definition.StageString() + ") in Pipelinefile")

    // execute all configured stages
    for _, stage := range definition.Stages {
        start := time.Now()
        log.Info(LOG_TAG, "executing stage \"" + stage.Name + "\"", "in image \"" + stage.Image + "\"")

        // execute the stage
        err := p.ExecuteStage(&stage)
        if err != nil {
            log.Error(LOG_TAG, "stage \"" + stage.Name + "\" failed:", err.Error())
            return err
        }

        log.Info(LOG_TAG, "stage \"" + stage.Name + "\" completed in", time.Since(start))
    }

    return nil
}

// Clones a copy of the source repository to the pipelines working dir.
func (p *Pipeline) Checkout() (error) {
    ret, err := p.Container(SCM_IMAGE, p.project.Repository)
    if err != nil {
        return err
    }

    if ret != 0 {
        return errors.New("error code:" + strconv.Itoa(ret))
    }

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
func (p *Pipeline) ExecuteStage(stage *pipelinefile.Stage) (error) {
    // construct the build steps command string
    steps := strings.Join(stage.Steps, " && ")
    steps = "/bin/sh -c '" + steps + "'"

    // default to configured image if neceasarry
    image := STAGE_IMAGE
    if len(stage.Image) > 0 {
        image = stage.Image
    }

    ret, err := p.Container(image, steps)
    if err != nil {
        return err
    }

    if ret != 0 {
        return errors.New("error code " + strconv.Itoa(ret))
    }

    return nil
}