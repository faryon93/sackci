package util
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
    "time"
    "errors"
)

// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

var (
    ErrCancel = errors.New("timer canceled")
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type Timer struct {
    time chan bool
    cancel chan bool
}


// --------------------------------------------------------------------------------------
//  constructors
// --------------------------------------------------------------------------------------

func NewTimer() (*Timer) {
    return &Timer{
        time:  make(chan bool),
        cancel: make(chan bool),
    }
}


// --------------------------------------------------------------------------------------
//  public members
// --------------------------------------------------------------------------------------

func (t *Timer) Wait(duration time.Duration) (error)  {
    go func() {
        time.Sleep(duration)
        t.time <- true
    }()

    select {
        case <- t.time:
            return nil

        case <- t.cancel:
            return ErrCancel
    }
}

func (t *Timer) Cancel()  {
    t.cancel <- true
}
