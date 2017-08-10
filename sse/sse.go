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
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

// Upgrades the connection to an SSE connection.
// A proper Content-Type is set and kee-alive is actived.
func Upgrade(w http.ResponseWriter) (error) {
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
func RegisterCloseHandler(w http.ResponseWriter, f func()) {
    notify := w.(http.CloseNotifier).CloseNotify()
    go func() {
        // wait until the close notification arrives
        <- notify
        f()
    }()
}

// Write a string to the SSE stream.
func Write(w http.ResponseWriter, message string) (error) {
    _, err := w.Write([]byte(message + "\n"))
    if err != nil {
        return err
    }

    // immediately send the content to the client
    w.(http.Flusher).Flush()

    return nil
}

// Writes an event to the SSE stream.
func WriteEvent(w http.ResponseWriter, event Event) (error) {
    // encode the event in json representation
    b, err := json.Marshal(event)
    if err != nil {
        return err
    }

    return Write(w, "event: " + event.Event() + "\n" +
                    "data: " + string(b) + "\n")
}
