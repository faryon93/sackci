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

// ----------------------------------------------------------------------------------
//  imports
// ----------------------------------------------------------------------------------

import (
    "os"

    log "github.com/sirupsen/logrus"
    "github.com/asdine/storm/q"
    "github.com/asdine/storm"
)


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

// Executes the garbage collection process for this project.
func (p *Project) RunGc(artifacts string) (error) {
    // we're finished if no gc is configured
    if p.KeepHistory < 1 {
        return nil
    }

    // delete all builds except the last n builds
    // which were set in configuration file
    var builds []Build
    query := Get().Select(q.Eq("Project", p.Id)).Reverse().Skip(p.KeepHistory)
    err := query.Find(&builds)
    if err == storm.ErrNotFound {
        return nil
    } else if err != nil {
        return err
    }

    // remove the artifact files from the fs
    for _, build := range builds {
        err := os.Remove(build.GetArtifactName(artifacts))
        if err != nil {
            log.Errorf("failed to remove artifact file:", err.Error())
            continue
        }
    }

    return query.Delete(new(Build))
}