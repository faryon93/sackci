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
    "io"
    "archive/tar"
    "bytes"
    "errors"
    "os"
    "compress/gzip"

    "github.com/fsouza/go-dockerclient"
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type tarFilesystem struct {
    buffer []byte
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (f *tarFilesystem) Write(p []byte) (int, error) {
    f.buffer = append(f.buffer, p...)
    return len(p), nil
}

func (f *tarFilesystem) Get(name string) ([]byte, error) {
    reader := tar.NewReader(bytes.NewReader(f.buffer))
    for {
        header, err := reader.Next()
        if err == io.EOF {
            return nil, errors.New("file not found")
        } else if err != nil {
            return nil, err
        }

        // we found the file we are searching for
        if header.Name == name {
            // read the content of the file
            buf := make([]byte, header.Size)
            reader.Read(buf)

            return buf, nil
        }
    }
}


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

func (a *Agent) ReadFile(container string, file string) ([]byte, error) {
    // construct the tarFilesystem
    fs := tarFilesystem{}

    // Download the file from the container
    err := a.docker.DownloadFromContainer(container, docker.DownloadFromContainerOptions{
        Path: a.Filepath(workdir, file),
        OutputStream: &fs,
    })
    if err != nil {
        return nil, err
    }

    return fs.Get(file)
}

// Saves the remote path to a local gzip compressed tar file.
func (a *Agent) SavePath(container string, path string, file string) (error) {
    // create the target file in local filesystem
    fh, err := os.Create(file)
    if err != nil {
        return err
    }
    defer fh.Close()

    // we need a gzip compressor
    zip := gzip.NewWriter(fh)
    defer zip.Close()

    return a.docker.DownloadFromContainer(container, docker.DownloadFromContainerOptions{
        Path: path,
        OutputStream: zip,
    })
}

// Joins parts of a filepath to a complete path.
func (a *Agent) Filepath(paths ...string) (string) {
    // the default linux implementation
    // if filepathes are constructed in another way
    // this would be the right place to implement
    p := ""
    for i, path := range paths {
        p += path

        if i < len(paths) - 1 {
            p += "/"
        }
    }

    return p
}