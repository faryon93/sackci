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
    "errors"
    "time"

    "github.com/gorilla/mux"
    "github.com/asdine/storm"
    "github.com/asdine/storm/q"
    "github.com/liip/sheriff"

    "github.com/faryon93/sackci/model"
)


// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    GROUP_QUERYALL = "queryall"
    GROUP_ALL      = "all"
    GROUP_ONE      = "one"

    HEADER_TIMESTAMP = "X-Timestamp"
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

// Queries all entries by one ore more fields.
// The entries are matched against the url routing parameters.
func QueryAll(router *mux.Router, path string, mod interface{}) (error) {
    modelType := reflect.SliceOf(reflect.TypeOf(mod))

    // make sure only slices and structs are registred
    if reflect.TypeOf(mod).Kind() != reflect.Struct {
        return errors.New("model must be struct")
    }

    handler := func(w http.ResponseWriter, r *http.Request) {
        // construct the query
        query, err := httpQuery(r)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotAcceptable)
            return
        }

        // query the database for all elements of the given model
        element := reflect.New(modelType)
        err = query.Find(element.Interface())
        if err == storm.ErrNotFound {
            // we want to return an empty slice if nothing
            // has been found
            element = reflect.MakeSlice(modelType, 0, 0)

        } else if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // send the filtered response
        filter(w, element.Interface(), GROUP_QUERYALL)
    }

    // register the corresponding routes in the router
    router.Methods("GET").Path(path).HandlerFunc(handler)

    return nil
}

// Queries just on entry by one ore more fields.
// The entries are matched against the url routing parameters.
func QueryOne(router *mux.Router, path string, mod interface{}) (error) {
    // make sure only structs are registred
    if reflect.TypeOf(mod).Kind() != reflect.Struct {
        return errors.New("model must be struct")
    }

    // fetch just one item of the model by its id
    one := func(w http.ResponseWriter, r *http.Request) {
        // construct the query
        query, err := httpQuery(r)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotAcceptable)
            return
        }

        // execute the query in database
        element := reflect.New(reflect.TypeOf(mod))
        err = query.First(element.Interface())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // send the filtered response
        filter(w, element.Interface(), GROUP_ONE)
    }

    // register the various handler functions
    router.Methods("GET").Path(path).HandlerFunc(one)

    return nil
}

// --------------------------------------------------------------------------------------
//  private functions
// --------------------------------------------------------------------------------------

func filter(w http.ResponseWriter, v interface{}, groups ...string) {
    // filter the output
    options := sheriff.Options{Groups: groups}
    filtered, err := sheriff.Marshal(&options, v)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // add a timestamp to each response
    w.Header().Set(HEADER_TIMESTAMP, strconv.FormatInt(time.Now().UnixNano(), 10))

    // send as json response
    Jsonify(w, filtered)
}

func httpQuery(r *http.Request) (storm.Query, error) {
    fields := mux.Vars(r)

    // construct the query expression
    matchers := []q.Matcher{}
    for field, value := range fields {
        // convert the url parameter to an integer
        fieldVal, err := strconv.Atoi(value)
        if err != nil {
            return nil, errors.New("faild to parse field \"" + field + "\"")
        }

        // All elements of the array are getting and'ed
        matchers = append(matchers, q.Eq(field, fieldVal))
    }

    // obtain the query which can be executed
    return model.Get().Select(matchers...), nil
}