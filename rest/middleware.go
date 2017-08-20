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
    "net/http"
    "reflect"
    "strconv"
    "strings"
    "time"

    "github.com/gorilla/mux"
    "github.com/asdine/storm/q"

    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/log"
    "github.com/asdine/storm"
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func Delete(router *mux.Router, path string, success func(*http.Request), mod interface{}, fields ...string) (error) {
    // make sure only slices and structs are registred
    if reflect.TypeOf(mod).Kind() != reflect.Struct {
        return ErrMustBeStruct
    }

    // fetch just one item of the model by its id
    deleteFn := func(w http.ResponseWriter, r *http.Request) {
        element := reflect.New(reflect.SliceOf(reflect.TypeOf(mod)))

        start := time.Now()

        // construct the query
        matchers := []q.Matcher{}
        for _, field := range fields {
            // parse the url parameters for the id
            fieldVal, err := strconv.Atoi(mux.Vars(r)[strings.ToLower(field)])
            if err != nil {
                http.Error(w, "invalid value for field \"" + field + "\"", http.StatusNotAcceptable)
                return
            }

            matchers = append(matchers, q.Eq(field, fieldVal))
        }

        // execute the query in database
        err := model.Get().Select(matchers...).Find(element.Interface())
        if err == storm.ErrNotFound {
        } else if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // delete all queried elements in the database
        element = element.Elem()
        for i := 0; i < element.Len(); i++ {
            err := model.Get().DeleteStruct(element.Index(i).Addr().Interface())
            if err != nil {
                log.Error("middleware", "failed to delete entry:", err.Error())
                continue
            }
        }

        log.Info("middleware",r.RequestURI, "took", time.Since(start), "to delete", element.Len(), "items")

        // call the success handler and return a json success
        success(r)
        Jsonify(w, struct{Success bool `json:"success"` }{true })

    }

    // register the various handler functions
    router.Methods(http.MethodDelete).Path(path).HandlerFunc(deleteFn)

    return nil
}
