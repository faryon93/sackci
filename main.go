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
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "context"
    "time"
    "runtime"

    "github.com/gorilla/mux"

    "github.com/faryon93/sackci/model"
    "github.com/faryon93/sackci/sse"
    "github.com/faryon93/sackci/ctx"
    "github.com/faryon93/sackci/rest"
    "github.com/faryon93/sackci/log"
)


// --------------------------------------------------------------------------------------
//  application entry
// --------------------------------------------------------------------------------------

func main() {
    log.Info("main", "starting sackci v0.1 #32h9d042v")

    // setup go environment
    runtime.GOMAXPROCS(runtime.NumCPU())

    // open database
    err := model.Open("E:/tmp/sackci")
    if err != nil {
        log.Error("bolt", "failed to open database:", err.Error())
        return
    }
    defer func() {
        model.Close()
        log.Info("bolt", "closed bolt database handle")
    }()
    log.Info("bolt", "successfully opened bolt database")

    // initialize the global application context
    ctx.Init()

    // create http server
    // and setup the routes with corresponding handler functions
    router := mux.NewRouter().StrictSlash(true)

    // register REST and SSE endpoints
    rest.Register(router, "/api/v1/project", model.Project{})
    rest.QueryAll(router, "/api/v1/project/{project}/env", "Project", model.Env{})
    sse.Register(router, "/api/v1/feed", ctx.Feed)

    // execute http server asynchronously
    srv := &http.Server{Addr: "127.0.0.1:8181", Handler: router}
    go func() {
        log.Info("http", "http server is listening on 0.0.0.0:8181")
        err := srv.ListenAndServe()
        if err != nil && err != http.ErrServerClosed {
            log.Error("http", "failed to serv http:", err.Error())
            return
        }

        log.Info("http", "http server is now closed")
    }()

    // wait for a signal to shutdown the application
    wait(os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    log.Info("main", "initiating application shutdown (SIGINT / SIGTERM)")

    // gracefully shutdown the http server
    httpCtx, _ := context.WithTimeout(context.Background(), 1*time.Second)
    srv.Shutdown(httpCtx)
}


// --------------------------------------------------------------------------------------
//  helper functions
// --------------------------------------------------------------------------------------

func wait(sig ...os.Signal) {
    signals := make(chan os.Signal)
    signal.Notify(signals, sig...)
    <- signals
}
