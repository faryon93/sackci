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
    "strings"
    "github.com/gorilla/mux"
    "github.com/asdine/storm"
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
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func All(router *mux.Router, path string, mod interface{}) (error) {
    // make sure only slices and structs are registred
    if reflect.TypeOf(mod).Kind() != reflect.Struct {
        return errors.New("model must be struct")
    }

    // fetches all items of the model
    all := func(w http.ResponseWriter, r *http.Request) {
        modelType := reflect.SliceOf(reflect.TypeOf(mod))
        element := reflect.New(modelType)

        // query the database for all elements of the given model
        err := model.Get().All(element.Interface())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // filter the output
        options := sheriff.Options{Groups: []string{GROUP_ALL}}
        filtered, err := sheriff.Marshal(&options, element.Interface())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // send as json response
        Jsonify(w, filtered)
    }

    // register the various handler functions
    router.Methods("GET").Path(path).HandlerFunc(all)

    return nil
}

// Registers a model in the REST interface router
func One(router *mux.Router, path string, mod interface{}) (error) {
    // make sure only slices and structs are registred
    if reflect.TypeOf(mod).Kind() != reflect.Struct {
        return errors.New("model must be struct")
    }

    // fetch just one item of the model by its id
    one := func(w http.ResponseWriter, r *http.Request) {
        element := reflect.New(reflect.TypeOf(mod))

        // parse the url parameters for the id
        id, err := strconv.Atoi(mux.Vars(r)["id"])
        if err != nil {
            http.Error(w, "invalid id", http.StatusNotAcceptable)
            return
        }

        // search the database for the field
        err = model.Get().One("Id", id, element.Interface())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // filter the output
        options := sheriff.Options{Groups: []string{GROUP_ONE}}
        filtered, err := sheriff.Marshal(&options, element.Interface())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // send as json response
        Jsonify(w, filtered)
    }

    // register the various handler functions
    router.Methods("GET").Path(path).HandlerFunc(one)

    return nil
}

// Queries a model by a field.
func QueryAll(router *mux.Router, path string, field string, mod interface{}) (error) {
    modelType := reflect.SliceOf(reflect.TypeOf(mod))

    // make sure only slices and structs are registred
    if reflect.TypeOf(mod).Kind() != reflect.Struct {
        return errors.New("model must be struct")
    }

    handler := func(w http.ResponseWriter, r *http.Request) {
        element := reflect.New(modelType)

        // parse the url parameters for the id
        id, err := strconv.Atoi(mux.Vars(r)[strings.ToLower(field)])
        if err != nil {
            http.Error(w, "invalid id", http.StatusNotAcceptable)
            return
        }

        // query the database for all elements of the given model
        err = model.Get().Find(field, id, element.Interface())
        if err == storm.ErrNotFound {
            // we want to return an empty slice if nothing
            // has been found
            element = reflect.MakeSlice(modelType, 0, 0)

        } else if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // filter the output
        options := sheriff.Options{Groups: []string{GROUP_QUERYALL}}
        filtered, err := sheriff.Marshal(&options, element.Interface())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // send as json response
        Jsonify(w, filtered)
    }

    // register the corresponding routes in the router
    router.Methods("GET").Path(path).HandlerFunc(handler)

    return nil
}
