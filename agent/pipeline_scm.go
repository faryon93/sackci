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
    "strconv"
    "errors"
    "encoding/json"

    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/util"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    // scm commands
    SCM_CLONE    = "clone"
    SCM_HEAD_REF = "head"

    // scm compare return codes
    SCM_RET_SUCCESS        = 0
    SCM_RET_INVALID_BRANCH = 1
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Clones a copy of the source repository to the pipelines working dir.
func (p *Pipeline) Clone() (*model.Commit, error) {
    args := SCM_CLONE + " " + p.project.GetRepository() + " " + p.project.Branch

    // the last line is the json representation of the commit
    lastLine := ""

    // start the special SCM container to clone the repository
    ret, err := p.Container(p.project.Scm, args, MOUNTPOINT, func(line string) {
        p.LogTerminal(STAGE_SCM_ID, util.MaskCredentials(line))
        lastLine = line
    })
    if err != nil {
        return nil,err
    }

    // make sure the clone process exited with a proper exit code
    if ret != 0 {
        return nil,errors.New("error code: " + strconv.Itoa(ret))
    }

    // parse the commit information
    var commit model.Commit
    return &commit, json.Unmarshal([]byte(lastLine), &commit)
}

// Gets the head reference of the given pipeline.
func (p *Pipeline) HeadRef() (string, error) {
    args := SCM_HEAD_REF + " " + p.project.GetRepository() + " " + p.project.Branch

    // the last line of the output will be the reference
    ref := ""

    // start the special SCM container to detect new changes
    ret, err := p.Container(p.project.Scm, args, WORKDIR, func(line string) {
        ref = line
    })
    if err != nil {
        return "", err
    }

    // some return codes have a special meaning...
    if ret == SCM_RET_SUCCESS {
        return ref, nil

    } else if ret == SCM_RET_INVALID_BRANCH {
        return "", ErrInvalidBranch

    } else {
        return "", errors.New("error code: " + strconv.Itoa(ret))
    }
}
