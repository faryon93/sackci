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
    "strings"
    "unicode"

    "github.com/fsouza/go-dockerclient"
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Executes a command on the build agent, while using the given volume and image.
func (a *Agent) Execute(vol string, image string, cmd string) (string, int, error) {
    // create the container in order to start it
    container, err := a.docker.CreateContainer(docker.CreateContainerOptions{
        Config: &docker.Config{
            Image: image,
            Cmd: removeOutter(splitQuoted(cmd), "'"),
            Tty: true,
            AttachStderr: true,
            AttachStdout: true,
            WorkingDir: workdir,
        },
        HostConfig: &docker.HostConfig{
            Binds: []string{vol + ":" + workdir},
        },
    })
    if err != nil {
        return container.ID, -1, err
    }

    // start the container
    err = a.docker.StartContainer(container.ID, nil)
    if err != nil {
        return container.ID, -1, nil
    }

    // gather all console outputs
    err = a.docker.Logs(docker.LogsOptions{
        Container: container.ID,
        OutputStream: &ConsoleOutput{},
        ErrorStream: &ConsoleOutput{},
        Stdout: true,
        Stderr: true,
        Follow: true,
        RawTerminal: true,
    })
    if err != nil {
        return container.ID, -1, err
    }

    // wait for the container to finish the command
    ret, err := a.docker.WaitContainer(container.ID)
    return container.ID, ret, err
}

func (a *Agent) RemoveContainer(id string) (error) {
    return a.docker.RemoveContainer(docker.RemoveContainerOptions{ID: id})
}


// ----------------------------------------------------------------------------------
//  private functions
// ----------------------------------------------------------------------------------

func splitQuoted(s string) ([]string) {
    lastQuote := rune(0)
    return strings.FieldsFunc(s, func(c rune) bool {
        switch {
        case c == lastQuote:
            lastQuote = rune(0)
            return false
        case lastQuote != rune(0):
            return false
        case unicode.In(c, unicode.Quotation_Mark):
            lastQuote = c
            return false
        default:
            return unicode.IsSpace(c)
        }
    })
}

func removeOutter(s []string, sep string) []string {
    for i := range s {
        if strings.HasPrefix(s[i], sep) {
            s[i] = s[i][1:len(s[i])]
        }

        if strings.HasSuffix(s[i], sep) {
            s[i] = s[i][0:len(s[i])-1]
        }
    }

    return s
}