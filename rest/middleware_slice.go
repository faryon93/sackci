package rest
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
    "strconv"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/faryon93/sackci/util"
)


// --------------------------------------------------------------------------------------
//  middleware handlers
// --------------------------------------------------------------------------------------

// Queries the first element of a slice matching the url fields.
func sliceQueryOne(w http.ResponseWriter, r *http.Request, slice reflect.Value) {
    elem := findInSlice(slice, mux.Vars(r))
    if elem == nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }

    filter(w, elem.Interface())
}

// Updates on struct inside the slice, matching the url fields.
func sliceUpdateOne(w http.ResponseWriter, r *http.Request, slice reflect.Value,
                    updates map[string]interface{}, success func(r *http.Request)) {
    // find the struct inside the given slice matching the url fields
    elem := findInSlice(slice, mux.Vars(r))
    if elem == nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }

    // range over the fields to update
    for field, value := range updates {
        structField := util.FieldByTag(*elem, "json", field)
        if !structField.IsValid() {
            continue
        }

        if v, ok := value.(int); ok {
            structField.SetInt(int64(v))
        } else if v, ok := value.(string); ok {
            structField.SetString(v)
        }
    }

    success(r)
}


// --------------------------------------------------------------------------------------
//  private helpers
// --------------------------------------------------------------------------------------

// Finds the element discribed by query in slice and returns it.
func findInSlice(slice reflect.Value, query map[string]string) (*reflect.Value) {
    // check all entries in the slice
    for i := 0; i < slice.Len(); i++ {
        elem := slice.Index(i)

        // return the first matching one
        if structMatches(elem.Interface(), query) {
            return &elem
        }
    }

    return nil
}

// Checks if a structs matches the given query fields.
func structMatches(o interface{}, query map[string]string) (bool) {
    matches := true
    for fieldName, fieldValue := range query {
        // check if the field is present in the struct.
        // if not we should not consider this struct as a matching one
        field := reflect.ValueOf(o).FieldByName(fieldName)
        if !field.IsValid() {
            matches = false
            break
        }

        // compare the integer
        if field.Type().Kind() == reflect.Int {
            intval, err := strconv.Atoi(fieldValue)
            if err != nil {
                matches = false
                break
            }

            if field.Interface() != intval {
                matches = false
                break
            }

            // compare string type
        } else if field.Type().Kind() == reflect.String {
            if field.String() != fieldValue {
                matches = false
                break
            }
        }
    }

    return matches
}
