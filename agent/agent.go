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
//  types
// ----------------------------------------------------------------------------------

type Agent struct {
    Name string `yaml:"name"`
    Endpoint string `yaml:"endpoint"`
    Concurrent int `yaml:"concurrent"`

    docker *docker.Client
    mutex sync.Mutex
    buildCount int
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (a *Agent) Free() {
    a.mutex.Lock()
    if a.buildCount > 0 {
        a.buildCount--
    }
    a.mutex.Unlock()
}