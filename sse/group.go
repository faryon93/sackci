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
    "sync"

    "github.com/faryon93/sackci/log"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type Group struct {
    Name string
    Mutex sync.Mutex
    Channels map[chan Event]bool
}

type Event interface {
    Event() string
}


// --------------------------------------------------------------------------------------
//  constructors
// --------------------------------------------------------------------------------------

func NewGroup(name string) (*Group) {
    return &Group{
        Name: name,
        Channels: make(map[chan Event]bool),
    }
}


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func Handler(group *Group, w http.ResponseWriter, r *http.Request) {
    // upgrade the connection to an SSE connection
    err := upgrade(w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Info(group.Name, "client", r.RemoteAddr, "opened sse connection")

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
            log.Error(group.Name, "failed to write feed to client:", err.Error())
            continue
        }
    }

    log.Info(group.Name, "client", r.RemoteAddr, "closed feed connection")
}


// --------------------------------------------------------------------------------------
//  public members
// --------------------------------------------------------------------------------------

// Publishs an event to the feed.
func (g *Group) Publish(action Event) {
    g.Mutex.Lock()
    defer g.Mutex.Unlock()

    // publish the message to all channels
    for ch := range g.Channels {
        ch <- action
    }
}

// Registers a new client to the feed.
func (g *Group) Register() (chan Event) {
    g.Mutex.Lock()
    defer g.Mutex.Unlock()

    // create the new channel and add to channel list
    ch := make(chan Event)
    g.Channels[ch] = true

    return ch
}

// Unregisters a client from the feed.
// As a side effect the clients channel gets closed.
func (g *Group) Unregister(ch chan Event){
    g.Mutex.Lock()
    defer g.Mutex.Unlock()

    // delete the channel from the list
    close(ch)
    delete(g.Channels, ch)
}
