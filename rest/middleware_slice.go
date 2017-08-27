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
)

// --------------------------------------------------------------------------------------
//  middleware handlers
// --------------------------------------------------------------------------------------

// Queries the first element of a slice matching the url fields.
func sliceQueryOne(w http.ResponseWriter, r *http.Request, slice interface{}) {
    query := mux.Vars(r)
    value := reflect.ValueOf(slice)

    // check all entries in the slice
    for i := 0; i < value.Len(); i++ {
        elem := value.Index(i).Interface()

        // return the first matching one
        if structMatches(elem, query) {
            filter(w, elem)
            return
        }
    }

    // none was found -> 404 handler
    http.Error(w, "not found", http.StatusNotFound)
}


// --------------------------------------------------------------------------------------
//  private helpers
// --------------------------------------------------------------------------------------

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
