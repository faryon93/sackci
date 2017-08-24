package agent
// dockertest
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

// ----------------------------------------------------------------------------------
//  import
// ----------------------------------------------------------------------------------

import (
    "sync"
    "time"

    "github.com/fsouza/go-dockerclient"
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type Agent struct {
    // from config file
    Name string `yaml:"name" json:"name"`
    Endpoint string `yaml:"endpoint" json:"endpoint"`
    Concurrent int `yaml:"concurrent" json:"concurrent"`

    // public runtime variables
    BuildCount int `json:"build_count"`

    // private runtime variables
    docker *docker.Client
    mutex sync.Mutex
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Returns if this agent is ready for a new build.
func (a *Agent) IsReady() (bool) {
    a.mutex.Lock()
    defer a.mutex.Unlock()

    // decrease the timeout for the ping
    a.docker.SetTimeout(DOCKER_TIMEOUT)
    defer a.docker.SetTimeout(2 * time.Hour)

    // some validity checks
    okay := true
    if a.Concurrent > 0 && a.BuildCount >= a.Concurrent {
        okay = false
    }

    // check connectivitiy to the build agent
    return a.docker.Ping() == nil && okay
}

func (a *Agent) Free() {
    a.mutex.Lock()
    if a.BuildCount > 0 {
        a.BuildCount--
    }
    a.mutex.Unlock()
}
