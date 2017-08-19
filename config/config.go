package config
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
    "io/ioutil"

    "gopkg.in/yaml.v2"

    "github.com/faryon93/sackci/agent"
    "github.com/faryon93/sackci/log"
    "github.com/faryon93/sackci/model"
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type Config struct {
    Listen string `yaml:"listen"`
    Artifacts string `yaml:"artifacts"`
    Database string `yaml:"database"`
    Agents []agent.Agent `yaml:"agents"`
    Projects []model.Project `yaml:"projects"`
}


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

func Load(path string) (*Config, error) {
    buf, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var conf Config
    err = yaml.Unmarshal(buf, &conf)
    if err != nil {
        return nil, err
    }

    // fill in the project ids
    // TODO: better way of id handling
    for i := 0; i < len(conf.Projects); i++ {
        conf.Projects[i].Id = i + 1
    }

    return &conf, nil
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Returns a project by its ID.
func (c *Config) GetProject(id int) (*model.Project) {
    // check array bounds
    index := id - 1
    if index < 0 || index >= len(c.Projects) {
        return nil
    }

    return &c.Projects[index]
}

// Prints some important information of the config
func (c *Config) Print() {
    for _, project := range c.Projects {
        log.Info("conf", "adding project", project.Name, "(" + project.Repository + ")")
    }

    log.Info("conf", "artifact storage location:", c.Artifacts)
}