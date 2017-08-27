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
    "time"
    "errors"

    "github.com/gorilla/mux"
    "github.com/liip/sheriff"
)


// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    // groups
    GROUP_QUERYALL = "queryall"
    GROUP_QUERYONE = "one"

    // custom headerfields
    HEADER_TIMESTAMP = "X-Timestamp"

    // HTTP GET query parameters for paging
    GET_QUERY_LIMIT = "limit"
    GET_QUERY_SKIP  = "skip"

    // flags
    QUERY_REVERSE = iota

    LOG_TAG = "middleware"
)

var (
    ErrInvalidRef = errors.New("ref must be struct or slice")
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

// Queries all entries by one ore more fields.
// The entries are matched against the url routing parameters.
func QueryAll(router *mux.Router, path string, ref interface{}, flags ...int) (error) {
    var fn func(w http.ResponseWriter, r *http.Request)

    // slice query
    kind := reflect.TypeOf(ref).Kind()
    if kind == reflect.Struct {
        fn = func(w http.ResponseWriter, r *http.Request) {
            stormQueryAll(w, r, reflect.SliceOf(reflect.TypeOf(ref)), flags...)
        }

    // unknown type to process
    } else {
        return ErrInvalidRef
    }

    // register the corresponding routes in the router
    router.Methods(http.MethodGet).Path(path).HandlerFunc(fn)

    return nil
}

// Queries just on entry by one ore more fields.
// The entries are matched against the url routing parameters.
func QueryOne(r *mux.Router, path string, ref interface{}) (error) {
    var fn func(w http.ResponseWriter, r *http.Request)

    // slice query
    kind := reflect.TypeOf(ref).Kind()
    if kind == reflect.Slice {
        fn = func(w http.ResponseWriter, r *http.Request) {
            sliceQueryOne(w, r, ref)
        }

    // struct query -> storm
    } else if kind == reflect.Struct {
        fn = func(w http.ResponseWriter, r *http.Request) {
            stormQueryOne(w, r, reflect.TypeOf(ref))
        }

    // unknown type to process
    } else {
        return ErrInvalidRef
    }

    // register the various handler functions
    r.Methods(http.MethodGet).Path(path).HandlerFunc(fn)

    return nil
}

// Deletes all url field matching entries.
// The success function is called afterwards.
func DeleteAll(r *mux.Router, path string, ref interface{}, success func(*http.Request)) (error) {
    var fn func(w http.ResponseWriter, r *http.Request)

    // slice query
    kind := reflect.TypeOf(ref).Kind()
    if kind == reflect.Struct {
        fn = func(w http.ResponseWriter, r *http.Request) {
            stormDeleteAll(w, r, reflect.TypeOf(ref), success)
        }

    // unknown type to process
    } else {
        return ErrInvalidRef
    }

    // register the various handler functions
    r.Methods(http.MethodDelete).Path(path).HandlerFunc(fn)

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
