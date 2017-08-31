package sse
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
    "encoding/json"
    "errors"
    "reflect"

    "github.com/gorilla/mux"

    "github.com/faryon93/sackci/log"
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

// Registers an SSE stream with the given router and Group.
func Register(router *mux.Router, path string, group *Group) {
    // the handler function
    fn := func(w http.ResponseWriter, r *http.Request) {
        handler(group, w, r)
    }

    // register the handler in the router
    router.Methods("GET").Path(path).HandlerFunc(fn)
}


// --------------------------------------------------------------------------------------
//  private functions
// --------------------------------------------------------------------------------------

func handler(group *Group, w http.ResponseWriter, r *http.Request) {
    // upgrade the connection to an SSE connection
    err := upgrade(w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // register for the update feed
    // and make sure the channel gets closed
    // when the client disconnects
    ch := group.Register()
    closeHandler(w, func() {
        group.Unregister(ch)
    })

    // write all new items form the feed to the client
    // this blocks until the client disconnects
    for action := range ch {
        err := writeEvent(w, action)
        if err != nil {
            log.Error(group.Name, "failed to write event to client:", err.Error())
            continue
        }
    }
}

// Upgrades the connection to an SSE connection.
// A proper Content-Type is set and kee-alive is actived.
func upgrade(w http.ResponseWriter) (error) {
    // if streaming is possible -> setup apropriat headers
    _, ok := w.(http.Flusher)
    if !ok {
        return errors.New("sse not supported")
    }

    // setup sse headers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Cache-Control", "no-cache")

    return nil
}

// Registers a handler function that is called when the client
// closes the connection.
func closeHandler(w http.ResponseWriter, f func()) {
    notify := w.(http.CloseNotifier).CloseNotify()
    go func() {
        // wait until the close notification arrives
        <- notify
        f()
    }()
}

// write a string to the SSE stream.
func write(w http.ResponseWriter, message string) (error) {
    _, err := w.Write([]byte(message + "\n"))
    if err != nil {
        return err
    }

    // immediately send the content to the client
    w.(http.Flusher).Flush()

    return nil
}

// Writes an event to the SSE stream.
func writeEvent(w http.ResponseWriter, event interface{}) (error) {
    // encode the event in json representation
    b, err := json.Marshal(event)
    if err != nil {
        return err
    }

    return write(w, "event: " + reflect.TypeOf(event).Name() + "\n" +
                    "data: " + string(b) + "\n")
}

