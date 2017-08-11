package scm
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
    "strings"
    "time"

    "gopkg.in/src-d/go-git.v4/plumbing/transport/client"
    "gopkg.in/src-d/go-git.v4/plumbing/transport"

    "github.com/faryon93/sackci/log"
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func GitWatch(url string, branch string) {
    session, err := git(url)
    if err != nil {
        log.Error("git", "failed to connect to git repo:", err.Error())
        return
    }

    // cyclically check for a new revision in the branch we are watching
    for {
        // get all the references
        ar, err := session.AdvertisedReferences()
        if err != nil {
            log.Error("git", "failed to get references:", err.Error())
            continue
        }

        refs, err := ar.AllReferences()
        if err != nil {
            log.Error("git", "failed to parse references:", err.Error())
            return
        }

        // find the reference of the correct branch
        found := false
        for _, ref := range refs {
            name := strings.Split(ref.Name().String(), "/")

            if ref.IsBranch() && name[2] == branch {
                log.Info("git", "branch reference:", ref.Hash().String())
                found = true
                break
            }
        }

        if !found {
            log.Error("git", "repository has no branch \"" + branch + "\"")
            return
        }

        time.Sleep(1 * time.Second)
    }
}

func git(url string) (transport.UploadPackSession, error) {
    endpoint, err := transport.NewEndpoint(url)
    if err != nil {
        return nil, err
    }

    clnt, err := client.NewClient(endpoint)
    if err != nil {
        return nil, err
    }

   return clnt.NewUploadPackSession(endpoint, nil)
 }
