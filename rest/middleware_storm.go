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
    "strconv"
    "errors"
    "reflect"
    "time"

    "github.com/asdine/storm"
    "github.com/gorilla/mux"
    "github.com/asdine/storm/q"

    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/log"
    "github.com/faryon93/sackci/util"
)

// --------------------------------------------------------------------------------------
//  middleware handlers
// --------------------------------------------------------------------------------------

func stormQueryOne(w http.ResponseWriter, r *http.Request, typ reflect.Type) {
    // construct the query
    query, err := stormQuery(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotAcceptable)
        return
    }

    // execute the query in database
    element := reflect.New(typ)
    err = query.First(element.Interface())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // send the filtered response
    filter(w, element.Interface(), GROUP_QUERYONE)
}

func stormQueryAll(w http.ResponseWriter, r *http.Request, typ reflect.Type, flags ...int) {
    // construct the query
    query, err := stormQuery(r, flags...)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotAcceptable)
        return
    }

    // query the database for all elements of the given model
    element := reflect.New(typ)
    err = query.Find(element.Interface())
    if err == storm.ErrNotFound {
        // we want to return an empty slice if nothing
        // has been found
        element = reflect.MakeSlice(typ, 0, 0)

    } else if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // send the filtered response
    filter(w, element.Interface(), GROUP_QUERYALL)
}

func stormDeleteAll(w http.ResponseWriter, r *http.Request, typ reflect.Type, success func(*http.Request)) {
    start := time.Now()

    query, err := stormQuery(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotAcceptable)
        return
    }

    // execute the query in database
    element := reflect.New(reflect.SliceOf(typ))
    err = query.Find(element.Interface())
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
            log.Error(LOG_TAG, "failed to delete entry:", err.Error())
            continue
        }
    }

    log.Info(LOG_TAG, "DELETE", r.RequestURI, "took",
        time.Since(start), "to delete", element.Len(), "items")

    // call the success handler and return a json success
    success(r)
    Jsonify(w, struct{Success bool `json:"success"` }{true })
}


// --------------------------------------------------------------------------------------
//  private helpers
// --------------------------------------------------------------------------------------

func stormQuery(r *http.Request, flags ...int) (storm.Query, error) {
    fields := mux.Vars(r)

    // construct the query expression
    matchers := []q.Matcher{}
    for field, value := range fields {
        // only check upper case url parameters
        if util.IsUpper(field[0]) {
            // convert the url parameter to an integer
            fieldVal, err := strconv.Atoi(value)
            if err != nil {
                return nil, errors.New("faild to parse field \"" + field + "\"")
            }

            // All elements of the array are getting and'ed
            matchers = append(matchers, q.Eq(field, fieldVal))
        }
    }
    query := model.Get().Select(matchers...)

    // add all modifieres to the query
    // according to the flags
    for _, flag := range flags {
        if flag == QUERY_REVERSE {
            query = query.Reverse()
        }
    }

    // apply the limit parameter
    limit := r.URL.Query().Get(GET_QUERY_LIMIT)
    if len(limit) > 0 {
        limitVal, err := strconv.Atoi(limit)
        if err != nil {
            return nil, errors.New("invalid limit value")
        }

        query = query.Limit(limitVal)
    }

    // apply the skip parameter
    skip := r.URL.Query().Get(GET_QUERY_SKIP)
    if len(skip) > 0 {
        skipVal, err := strconv.Atoi(skip)
        if err != nil {
            return nil, errors.New("invalid skip value")
        }

        query = query.Skip(skipVal)
    }

    // obtain the query which can be executed
    return query, nil
}
