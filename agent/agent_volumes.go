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
    "github.com/fsouza/go-dockerclient"
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Create a new volume.
func (a *Agent) CreateVolume() (string, error) {
    vol, err := a.docker.CreateVolume(docker.CreateVolumeOptions{})
    return vol.Name, err
}

// Deletes the given volume by its name.
func (a *Agent) RemoveVolume(name string) (error) {
    return a.docker.RemoveVolume(name)
}
