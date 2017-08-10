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
    "encoding/binary"
    "github.com/boltdb/bolt"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type Insertable interface {
    SetId(uint64 uint64)
}


// --------------------------------------------------------------------------------------
//  private functions
// --------------------------------------------------------------------------------------

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}

func btoi(b []byte) uint64 {
    return binary.BigEndian.Uint64(b)
}

// makes sure that the buckets will be created
func assertBuckets(buckets []string) (error) {
    return db.Update(func(tx *bolt.Tx) error {
        for _, bucket := range buckets {
            _, err := tx.CreateBucketIfNotExists([]byte(bucket))
            if err != nil {
                return err
            }
        }

        return nil
    })
}