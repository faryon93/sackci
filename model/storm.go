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
    "github.com/asdine/storm"
    "github.com/asdine/storm/codec/gob"
    "github.com/faryon93/sackci/log"
)


// --------------------------------------------------------------------------------------
//  global variables
// --------------------------------------------------------------------------------------

// tiedot database handle
var db *storm.DB


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

// Opens the bolt database.
func Open(path string) (error) {
    var err error

    // open the database directory an keep the handle for later
    db, err = storm.Open(path, storm.Codec(gob.Codec))
    if err != nil {
        return err
    }

    return nil
}

// Closes the bolt database.
func Close() (error) {
    log.Info("bolt", "closed bolt database handle")
    return db.Close()
}

// Returns the database handle.
func Get() (*storm.DB) {
    return db
}
