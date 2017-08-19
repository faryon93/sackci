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

    "github.com/fsouza/go-dockerclient"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    STATUS_READY = "ready"
    STATUS_UNREACHABLE = "unreachable"
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
    Status string `json:"status"`

    // private runtime variables
    docker *docker.Client
    mutex sync.Mutex
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Returns if this agent is ready for a new build.
func (a *Agent) IsReady() (bool) {
    return a.Status == STATUS_READY
}

func (a *Agent) Free() {
    a.mutex.Lock()
    if a.BuildCount > 0 {
        a.BuildCount--
    }
    a.mutex.Unlock()
}