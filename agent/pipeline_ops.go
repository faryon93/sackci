package agent
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

// ----------------------------------------------------------------------------------
//  imports
// ----------------------------------------------------------------------------------

import (
    "github.com/faryon93/sackci/log"
    "github.com/faryon93/sackci/util"
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Executes a command in the given image on this pipeline.
func (p *Pipeline) Container(image string, cmd string, stdio func(string)) (int, error) {
    container, ret, err := p.Agent.Execute(p.Volume, image, cmd, stdio)
    if container != "" {
        p.mutex.Lock()
        p.Containers = append(p.Containers, container)
        p.mutex.Unlock()
    }

    return ret, err
}

// Read a file from the pipeline. At least on container has to be executed.
func (p *Pipeline) ReadFile(path string) ([]byte, error) {
    // to read a file we need at least on container to
    // access the volume
    if len(p.Containers) <= 0 {
        return nil, ErrNoContainer
    }

    return p.Agent.ReadFile(p.Containers[0], path)
}

// Saves a whole path from the pipeline as gzipped tar archive.
// At least on container has to be executed.
func (p *Pipeline) SavePath(path string, local string) (error) {
    // to read a file we need at least on container to
    // access the volume
    if len(p.Containers) <= 0 {
        return ErrNoContainer
    }

    return p.Agent.SavePath(p.Containers[0], path, local)
}

// Destroys the whole pipeline
func (p *Pipeline) Destroy() {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    // remove all containers
    for _, container := range p.Containers {
        err := p.Agent.RemoveContainer(container)
        if err != nil {
            log.Error(LOG_TAG, "failed to remove container:", err.Error())
            continue
        }
        log.Info(LOG_TAG, "removed container", util.ShortHash(container))
    }

    // destroy the volume
    err := p.Agent.RemoveVolume(p.Volume)
    if err != nil {
        log.Error(LOG_TAG, "failed to remove volume:", err.Error())
    }
    log.Info(LOG_TAG, "removed volume", util.ShortHash(p.Volume))

    // free the agent
    p.Agent.Free()

    // close the event stream
    close(p.Events)
}