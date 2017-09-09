package config
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
    "io/ioutil"
    "path/filepath"

    "gopkg.in/yaml.v2"

    "github.com/faryon93/sackci/agent"
    "github.com/faryon93/sackci/log"
    "github.com/faryon93/sackci/model"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    DATABASE = "meta.db"
    ARTIFACTS = "artifacts"
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type Config struct {
    Listen string `yaml:"listen,omitempty"`
    TlsKey string `yaml:"tlskey,omitempty"`
    TlsCert string `yaml:"tlscert,omitempty"`
    DataDir string `yaml:"datadir,omitempty"`
    Agents []agent.Agent `yaml:"agents,omitempty"`
    Projects []model.Project `yaml:"projects,omitempty"`

    // path of the config file this
    // instance was loaded from
    loadPath string
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
    conf.loadPath = path

    // fill in the project ids
    for i := 0; i < len(conf.Projects); i++ {
        // TODO: better way of id handling
        conf.Projects[i].Id = i + 1

        // an empty trigger means manual triggering
        if conf.Projects[i].Trigger == "" {
            conf.Projects[i].Trigger = model.TRIGGER_MANUAL
        }
    }

    return &conf, nil
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Saves the configuration file to the filesystem.
func (c *Config) Save() (error) {
    bytes, err := yaml.Marshal(c)
    if err != nil {
        return err
    }

    return ioutil.WriteFile(c.loadPath, bytes, 0644)
}

// Returns a project by its ID.
func (c *Config) GetProject(id int) (*model.Project) {
    // check array bounds
    index := id - 1
    if index < 0 || index >= len(c.Projects) {
        return nil
    }

    return &c.Projects[index]
}

// Returns the artifacts directory
func (c *Config) GetArtifactsDir() string {
    return filepath.Join(c.DataDir, ARTIFACTS)
}

// Returns the path to the database file
func (c *Config) GetDatabaseFile() string {
    return filepath.Join(c.DataDir, DATABASE)
}

// Prints some important information of the config
func (c *Config) Print() {
    for _, project := range c.Projects {
        log.Info("conf", "adding project", project.Name, "(" + project.Repository + ")")
    }

    log.Info("conf", "artifact storage location:", c.GetArtifactsDir())
}