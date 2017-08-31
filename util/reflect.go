package util
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
    "reflect"
)


// --------------------------------------------------------------------------------------
//  public members
// --------------------------------------------------------------------------------------

func FieldByTag(ref reflect.Value, tag string, tagValue string) (reflect.Value) {
    if ref.Type().Kind() != reflect.Struct {
        return reflect.Value{}
    }

    typ := ref.Type()
    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        if field.Tag.Get(tag) == tagValue {
            return ref.Field(i)
        }
    }

    return reflect.Value{}
}
