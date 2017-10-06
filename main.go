package main
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

// --------------------------------------------------------------------------------------
//  imports
// --------------------------------------------------------------------------------------

import (
    "os"
    "os/signal"
    "syscall"
    "time"
    "runtime"
    "math/rand"
    "flag"

    "github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"

    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/ctx"
    "github.com/faryon93/sackci/config"
    "github.com/faryon93/sackci/agent"
    "github.com/faryon93/sackci/scm"
)


// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    DEFAULT_CONFIG = "sackci.conf"
    HTTP_SHUTDOWN_TIMEOUT = 1 * time.Second
)


// --------------------------------------------------------------------------------------
//  global variables
// --------------------------------------------------------------------------------------

var (
    // command line options
    configPath string
    purge bool
    color bool
    timestamp bool
)


// --------------------------------------------------------------------------------------
//  application entry
// --------------------------------------------------------------------------------------

func main() {
    // parse command line arguments
    flag.StringVar(&configPath, "conf", DEFAULT_CONFIG, "path to config file")
    flag.BoolVar(&purge, "purge", false, "purge unreferenced metadata")
    flag.BoolVar(&color, "color", false, "force logging with colors")
    flag.BoolVar(&timestamp, "timestamp", false, "prepend full timestamps")
    flag.Parse()

    // setup the logging
    formater := log.TextFormatter{ForceColors: color, FullTimestamp: timestamp}
    log.SetFormatter(&formater)
    log.SetOutput(os.Stdout)

    // begin with the application startup
    log.Infoln("starting", GetAppVersion())

    // setup go environment
    runtime.GOMAXPROCS(runtime.NumCPU())
    rand.Seed(time.Now().Unix())

    // load the configuration file
    conf, err := config.Load(configPath)
    if err != nil {
        log.Errorln("failed to load config:", err.Error())
        return
    }
    ctx.Conf = conf

    // open database
    err = model.Open(ctx.Conf.GetDatabaseFile())
    if err != nil {
        log.Errorln("failed to open database:", err.Error())
        return
    }
    defer model.Close()
    log.Infoln("successfully opened bolt database")

    // initialize the global application context
    // purge all unreferenced metadata if user wants to
    ctx.Conf.Setup()
    if purge {
        model.PurgeProjects(ctx.Conf.Projects)
    }

    // setup agents and SCM polling
    agent.Add(ctx.Conf.Agents...)
    scm.Setup()
    defer scm.Destroy()

    // create http server
    // and setup the routes with corresponding handler functions
    router := mux.NewRouter().StrictSlash(true)
    routes(router)

    // setup http and https servers and make sure the are
    // shutdown gracefully
    srv := SetupHttpEndpoint(conf, router)
    defer ShutdownHttp(srv, HTTP_SHUTDOWN_TIMEOUT)
    if conf.IsHttpsEnabled() {
        srv := SetupHttpsEndpoint(conf, router)
        defer ShutdownHttp(srv, HTTP_SHUTDOWN_TIMEOUT)
    }

    log.Infoln("everything is now up and running, ready to build!")

    // wait for a signal to shutdown the application
    wait(os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    log.Infoln("initiating application shutdown (SIGINT / SIGTERM)")
}


// --------------------------------------------------------------------------------------
//  helper functions
// --------------------------------------------------------------------------------------

func wait(sig ...os.Signal) {
    signals := make(chan os.Signal)
    signal.Notify(signals, sig...)
    <- signals
}
