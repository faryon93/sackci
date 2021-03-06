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
    "time"

    "gopkg.in/yaml.v2"
    log "github.com/sirupsen/logrus"

    "github.com/faryon93/sackci/agent"
    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/util"
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
    HttpListen  string            `yaml:"http_listen,omitempty"`
    HttpsListen string            `yaml:"https_listen,omitempty"`
    HttpsKey    string            `yaml:"https_key,omitempty"`
    HttpsCert   string            `yaml:"https_cert,omitempty"`
    DataDir     string            `yaml:"datadir,omitempty"`
    Users       UserList          `yaml:"users,omitempty"`
    Agents      []agent.Agent     `yaml:"agents,omitempty"`
    Projects    []*model.Project   `yaml:"projects,omitempty"`

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

    return &conf, nil
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (c *Config) Setup() {
    start := time.Now()

    // each project has to be completed for runtime
    for _, project := range c.Projects {
        err := project.IsMissconfigured()
        if err != nil {
            log.Errorln("ignoring project", project.Name + ":", err.Error())
            continue
        }

        // assign a project id wich is fixed
        err = project.AssignId()
        if err != nil {
            log.Infoln("ignoring project \"" + project.Name + "\":", err.Error())
            project.Id = -1
            continue
        }

        // check the project integrity
        project.CheckIntegrity()

        // an empty trigger means manual triggering
        if project.Trigger == "" {
            project.Trigger = model.TRIGGER_MANUAL
        }

        // run the garbage collection
        err = project.RunGc(c.GetArtifactsDir())
        if err != nil {
            log.Errorf("gc for project \"%s\" failed: %s", project.Name, err.Error())
        }

        err = project.CreateLuaVm()
        if err != nil {
            log.Errorf("failed to create lua vm:", err.Error())
        }

        // everything was fine -> we want to keep this project in our list
        log.Infoln("adding project", project.Name, "(" + util.MaskCredentials(project.Repository) + ")")
    }

    // save config file to disk -> a hash might have been inserted
    // TODO: only save when necessary
    err := c.Save()
    if err != nil {
        log.Errorln("failed to save conf file:", err.Error())
    }

    log.Infoln( "project integrity check took", time.Since(start))

    // display a warning when authentication is disabled
    if !c.IsAuthEnabled() {
        log.Warn("AUTHENTICATION IS DISABLED! SOMEONE MAY STEAL YOUR SENSITIVE DATA!")
    }
}

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
    // serach for project with given id field
    for _, project := range c.Projects {
        if project.IsValid() && project.Id == id {
            return project
        }
    }

    return nil
}

// Returns the artifacts directory.
func (c *Config) GetArtifactsDir() string {
    return filepath.Join(c.DataDir, ARTIFACTS)
}

// Returns the path to the database file.
func (c *Config) GetDatabaseFile() string {
    return filepath.Join(c.DataDir, DATABASE)
}

// Returns true if http server is enabled.
func (c *Config) IsHttpEnabled() bool {
    return c.HttpListen != ""
}

// Returns true if ssl encrypted http(s) server is enabled.
func (c *Config) IsHttpsEnabled() bool {
    return !util.StrEmpty(c.HttpsListen, c.HttpsCert, c.HttpsKey)
}

// Returns true if user authentication is enabled.
func (c *Config) IsAuthEnabled() (bool) {
    return len(c.Users) > 0
}
