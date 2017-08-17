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
//  imports
// ----------------------------------------------------------------------------------

import (
    "sync"
    "math/rand"
    "time"

    "github.com/fsouza/go-dockerclient"

    "github.com/faryon93/sackci/log"
)

// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    DOCKER_TIMEOUT = 500 * time.Millisecond
)


// ----------------------------------------------------------------------------------
//  global variables
// ----------------------------------------------------------------------------------

var (
    poolMutex = sync.Mutex{}
    pool      = []Agent{}
    workdir   = "/"
)


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

// Adds one or multiple agents to the build agent pool.
func Add(agents ...Agent) {
    poolMutex.Lock()
    defer poolMutex.Unlock()

    // create the docker client connection
    for _, agent := range agents {
        // construct docker client
        // TODO: tls authentication for remote api
        client, err := docker.NewClient(agent.Endpoint)
        if err != nil {
            log.Error("agent", "failed to create docker client:", err.Error())
            continue
        }
        client.SetTimeout(DOCKER_TIMEOUT)

        // check connectivitiy to the build agent
        err = client.Ping()
        if err != nil {
            log.Error("agent", err.Error())
            return
        }

        // populate the agent with necesarry runtime fields
        agent.BuildCount = 0
        agent.docker = client

        // add the agent to the agent pool
        pool = append(pool, agent)
        log.Error("agent", "adding build agent", agent.Name, "(" + agent.Endpoint + ")")
    }
}

// Return a random build agent.
func Allocate() (*Agent) {
    if len(pool) < 1 {
        return nil
    }

    agent := &pool[rand.Intn(len(pool))]

    // increase the counters
    agent.mutex.Lock()
    agent.BuildCount++
    agent.mutex.Unlock()

    return agent
}

// Sets the workdir for alle containers.
func SetWorkdir(dir string) {
    workdir = dir
}