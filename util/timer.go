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
    ErrTimerCancel = errors.New("timer canceled")
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type CycleTimer struct {
    StartTime time.Time

    timeout   chan bool
    cancel    chan bool
    cycleTime time.Duration
    fn        func(t *CycleTimer)
    done      func()
}


// --------------------------------------------------------------------------------------
//  constructors
// --------------------------------------------------------------------------------------

func NewTimer(cycleTime time.Duration, fn func(t *CycleTimer), done func()) (*CycleTimer) {
    if fn == nil {
        return nil
    }

    return &CycleTimer{
        timeout:    make(chan bool),
        cancel:     make(chan bool),
        cycleTime:  cycleTime,
        fn:         fn,
        done:       done,
    }
}


// --------------------------------------------------------------------------------------
//  public members
// --------------------------------------------------------------------------------------

// Wait for the given amount of time.
func (t *CycleTimer) Wait(duration time.Duration) (error)  {
    go func() {
        time.Sleep(duration)
        t.timeout <- true
    }()

    select {
        case <- t.timeout:
            return nil

        case <- t.cancel:
            return ErrTimerCancel
    }
}

// Canceling the waiting of a cycle.
func (t *CycleTimer) Cancel()  {
    t.cancel <- true
}

func (t *CycleTimer) Start() {
    go t.run()
}


// --------------------------------------------------------------------------------------
//  private members
// --------------------------------------------------------------------------------------

func (t *CycleTimer) run() {
    // if the cycle loop exits, run the done funtion if necessary
    if t.done != nil {
        defer t.done()
    }

    for {
        // sleep for the cycle time
        err := t.Wait(t.cycleTime)
        if err != nil {
            return
        }

        // run th supplied function
        t.StartTime = time.Now()
        t.fn(t)
    }
}
