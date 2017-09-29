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
    "os"
    "path/filepath"

    log "github.com/sirupsen/logrus"
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Executes a command in the given image on this pipeline.
func (p *Pipeline) Container(image string, cmd string, workdir string, stdio func(string)) (int, error) {
    // create the container
    container, err := p.Agent.CreateContainer(p.Volume, image, cmd, p.getEnv(), workdir)
    if container != "" {
        p.addContainer(container)
    }
    if err != nil {
        return -1, err
    }

    // copy the key file to the container
    buf, err := p.project.GetPrivateKey()
    if err == nil {
        // if there was no problem loaded the private key
        // it should be uploaded to the container
        // just write the privatekey to the container if there
        // is a privatekey configured
        if len(buf) > 0 {
            err = p.Agent.WriteFile(container, KEY_PATH, buf, KEY_PERMISSIONS)
            if err != nil {
                return -1, err
            }
        }

    // there was an error loading the private key
    } else {
        log.Errorln("failed to get private key:", err.Error())
    }

    // execute the container
    return p.Agent.StartContainer(container, stdio)
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

    err := os.MkdirAll(filepath.Dir(local), os.ModePerm)
    if err != nil {
        return err
    }

    return p.Agent.SavePath(p.Containers[0], path, local)
}
