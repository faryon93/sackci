# sackci
Simple as *uck continuous integration server

![Screenshot of webinterface](doc/webinterface.png)

## Features
- agent-less design (just a docker host is needed)
- pipeline configuration stored in project repository
- configuration in text format
- custom scm via special docker container
- realtime web interface
- REST interface
- single binary, which contains all necessary data

## Getting started
The version contained in the master branch is packaged in a docker image and published on docker hub.
A valid config file must be mounted to the container by the user.
```
docker run --name sackci -t -i \
           -p 443:443 \
           -p 80:80 \
           -v /var/lib/sackci:/sackci \
           -v /var/run/docker.sock:/var/run/docker.sock \
           sackci/server
```

## Architecture
The sackci server is the central component which orchestrates the builds on all "agents". On the build agents there is no software required, except a running docker daemon with exposed management api to the network.
The server periodically polls for changes changes in the source repository of each project. If new changes are detected the build process is executed on one of the configured docker hosts.
All build stages are executed in seperate containers, so there is no need to install all the tools on the host machine. Just pick the right docker image for your build stage.

## SCM Integration
It is possible to integrate any kind of SCM into sackci. Therefore a special docker container is needed.
The sackci server communicates with the scm container via command line arguments and return values.
We provide an scm integration container for the following scm systems:

- [Git](https://github.com/faryon93/sackci-git)

## Required Tools
To build the webfrontend some tools are required to process less files, minify the content and embed all assets into the application.
All files are stored in *assets/fs.go*, which should be up to date at any time. In order to generate a new *fs.go* file a
simple `go generate assets/fs_util.go` should be enough.

- **minify**: https://github.com/tdewolff/minify
- **esc**: https://github.com/mjibson/esc
- **lessc**: https://lesscss.org/

## TODO
- secret variables in Pipelinefile
- graphing of code metrics
- hot reload of configuration
- project inheritance 

## Notice
This project is far from beeing finished and should not be used in production. Feel free to contribute.
