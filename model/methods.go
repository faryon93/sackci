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
    "errors"
    "encoding/json"

    "github.com/boltdb/bolt"
)


// --------------------------------------------------------------------------------------
//  private functions
// --------------------------------------------------------------------------------------

func Insert(bucket string, v Insertable) (uint64, error) {
    id := uint64(0)
    return id, db.Update(func (tx *bolt.Tx) error {
        // Get the bucket handle
        b := tx.Bucket([]byte(bucket))
        if b == nil {
            return errors.New("Insert: bucket \"" + bucket + "\" does not exist")
        }

        // Generate a unique id.
        // This returns an error only if the Tx is closed or not writeable.
        // That can't happen in an Update() call so I ignore the error check.
        id, _ = b.NextSequence()
        v.SetId(id)

        // Marshal data into bytes
        buf, err := json.Marshal(v)
        if err != nil {
            return err
        }

        return b.Put(itob(id), buf)
    })
}

func InsertId(bucket string, id uint64, v interface{}) (error) {
    return db.Update(func (tx *bolt.Tx) error {
        // Get the bucket handle
        b := tx.Bucket([]byte(bucket))
        if b == nil {
            return errors.New("InsertId: bucket \"" + bucket + "\" does not exist")
        }

        // Marshal data into bytes
        buf, err := json.Marshal(v)
        if err != nil {
            return err
        }

        return b.Put(itob(id), buf)
    })
}

func Get(bucket string, id uint64, v interface{}) (error) {
    return db.View(func (tx *bolt.Tx) error {
        // Get the bucket handle
        b := tx.Bucket([]byte(bucket))
        if b == nil {
            return errors.New("get: bucket \"" + bucket + "\" does not exist")
        }

        buf := b.Get(itob(id))
        if buf == nil {
            return errors.New("key does not exist")
        }

        // unmarshall the json
        return json.Unmarshal(buf, v)
    })
}
