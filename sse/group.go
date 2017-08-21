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
    "sync"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type Group struct {
    Name string
    Mutex sync.Mutex
    Channels map[chan interface{}]bool
}


// --------------------------------------------------------------------------------------
//  constructors
// --------------------------------------------------------------------------------------

func NewGroup(name string) (*Group) {
    return &Group{
        Name: name,
        Channels: make(map[chan interface{}]bool),
    }
}


// --------------------------------------------------------------------------------------
//  public members
// --------------------------------------------------------------------------------------

// Publishs an event to the feed.
func (g *Group) Publish(event interface{}) {
    g.Mutex.Lock()
    defer g.Mutex.Unlock()

    // publish the message to all channels
    for ch := range g.Channels {
        ch <- event
    }
}

// Registers a new client to the feed.
func (g *Group) Register() (chan interface{}) {
    g.Mutex.Lock()
    defer g.Mutex.Unlock()

    // create the new channel and add to channel list
    ch := make(chan interface{})
    g.Channels[ch] = true

    return ch
}

// Unregisters a client from the feed.
// As a side effect the clients channel gets closed.
func (g *Group) Unregister(ch chan interface{}){
    g.Mutex.Lock()
    defer g.Mutex.Unlock()

    // delete the channel from the list
    close(ch)
    delete(g.Channels, ch)
}
