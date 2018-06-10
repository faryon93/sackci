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
    "strings"

    "github.com/fsouza/go-dockerclient"
    "github.com/kballard/go-shellquote"
)

// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

const (
    DOCKER_SOCKET = "/var/run/docker.sock"
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func(a *Agent) CreateContainer(vol string, image string, cmd string, env []string, workdir string) (string, error) {
    // setup the mount points use in container
    volumes := []string {
        vol + ":" + MOUNTPOINT,
    }

    // when using the special docker image, a socket to the
    // host docker daemon is needed as well
    if strings.Contains(image, "docker") {
        volumes = append(volumes, DOCKER_SOCKET + ":" + DOCKER_SOCKET)
    }

    // split the command like the shell would do it
    command, err := shellquote.Split(cmd)
    if err != nil {
        return "", err
    }

    // create the container in order to start it
    container, err := a.docker.CreateContainer(docker.CreateContainerOptions{
        Config: &docker.Config{
            Image: image,
            Cmd: command,
            Tty: true,
            AttachStderr: true,
            AttachStdout: true,
            WorkingDir: workdir,
            Env: env,
        },
        HostConfig: &docker.HostConfig{Binds: volumes},
    })
    if err != nil && container == nil {
        return "", err
    }

    return container.ID, err
}

// Executes a command on the build agent, while using the given volume and image.
func (a *Agent) StartContainer(container string, stdio func(string)) (int, error) {
    // start the container
    err := a.docker.StartContainer(container, nil)
    if err != nil {
        return -1, err
    }

    // gather all console outputs
    err = a.docker.Logs(docker.LogsOptions{
        Container: container,
        OutputStream: &ConsoleOutput{Callback: stdio},
        ErrorStream: &ConsoleOutput{Callback: stdio},
        Stdout: true,
        Stderr: true,
        Follow: true,
        RawTerminal: true,
    })
    if err != nil {
        return -1, err
    }

    // wait for the container to finish the command
    ret, err := a.docker.WaitContainer(container)
    return ret, err
}

// Removes a container with the given ID from the build agent.
func (a *Agent) RemoveContainer(id string) (error) {
    return a.docker.RemoveContainer(docker.RemoveContainerOptions{ID: id})
}
