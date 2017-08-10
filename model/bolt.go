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
    "log"
    "github.com/boltdb/bolt"
)


// --------------------------------------------------------------------------------------
//  global variables
// --------------------------------------------------------------------------------------

// tiedot database handle
var db *bolt.DB


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

// Opens the bolt database.
func Open(path string) (error) {
    var err error

    // open the database directory an keep the handle for later
    db, err = bolt.Open(path, 0600, nil)
    if err != nil {
        log.Println("failed to open database:", err.Error())
        return err
    }

    // make sure that some buckets exist
    err = assertBuckets([]string{PROJECTS_BUCKET, ENV_BUCKET, BUILDS_BUCKET, BUILDS_INDEX_BUCKET})
    if err != nil {
        log.Println("failed to assert buckets:", err.Error())
        return err
    }

    return nil
}

// Closes the bolt database.
func Close() (error) {
    return db.Close()
}






