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
    "github.com/boltdb/bolt"

    "encoding/json"
    "time"
)

// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    PROJECTS_BUCKET = "projects"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type Project struct {
    Id          uint64      `json:"id"`
    Name        string      `json:"name"`
    BuildStatus string      `json:"status"`
    BuildId     int         `json:"build"`
    Time        time.Time   `json:"execution_time"`
    Duration    string      `json:"duration"`
}


// --------------------------------------------------------------------------------------
//  public members
// --------------------------------------------------------------------------------------

func (p *Project) SetId(id uint64) {
    p.Id = id
}

func (p *Project) GetHistory() ([]Build, error) {
    return GetProjectHistory(p.Id)
}


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func QueryProjects() ([]Project, error) {
    projects := make([]Project, 0)

    return projects, db.View(func(tx *bolt.Tx) error {
        // Assume bucket exists and has keys
        b := tx.Bucket([]byte(PROJECTS_BUCKET))

        c := b.Cursor()
        for k, v := c.First(); k != nil; k, v = c.Next() {
            // parse the json string
            var project Project
            err := json.Unmarshal(v, &project)
            if err != nil {
                return err
            }

            projects = append(projects, project)
        }

        return nil
    })
}