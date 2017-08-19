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
    "github.com/faryon93/sackci/model"
    "encoding/json"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    // scm commands
    SCM_CLONE = "clone"
    SCM_COMPARE = "compare"

    // scm compare return codes
    SCM_RET_NEW_REF = 0
    SCM_RET_NO_CHANGES = 1
    SCM_RET_INVALID_BRANCH = 2
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Clones a copy of the source repository to the pipelines working dir.
func (p *Pipeline) Clone() (*model.Commit, error) {
    args := SCM_CLONE + " " + p.project.Repository + " " + p.project.Branch

    // the last line is the json representation of the commit
    lastLine := ""

    // start the special SCM container to clone the repository
    ret, err := p.Container(SCM_IMAGE, args, func(line string) {
        p.Events.ConsoleLog(STAGE_SCM_ID, line)
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

// Checks if a new commit is available.
func (p *Pipeline) Compare(old string) (bool, string, error) {
    args := SCM_COMPARE + " " + p.project.Repository + " " + p.project.Branch + " " + old

    // the last line of the output will be the reference
    ref := ""

    // start the special SCM container to detect new changes
    ret, err := p.Container(SCM_IMAGE, args, func(line string) {
        ref = line
    })
    if err != nil {
        return false, "", err
    }

    // some return codes have a special meaning...
    if ret == SCM_RET_NEW_REF {
        return true, ref, nil

    } else if ret == SCM_RET_NO_CHANGES {
        return false, old, nil

    } else if ret == SCM_RET_INVALID_BRANCH {
        return false, "", ErrInvalidBranch

    } else {
        return false, "", errors.New("error code: " + strconv.Itoa(ret))
    }
}