package model
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
    "time"

    "github.com/asdine/storm/q"
    "github.com/asdine/storm"
    log "github.com/sirupsen/logrus"
)


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

func PurgeProjects(projects []Project) {
    start := time.Now()

    // construct the query for all invalid projects
    // in order to delete
    invalidProjects := make([]q.Matcher, len(projects))
    for i, project := range projects {
        invalidProjects[i] = q.Not(q.Eq("Id", project.Id))
    }

    query := Get().Select(invalidProjects...)
    count, err := query.Count(new(ProjectMapping))
    if err != nil {
        log.Errorln("failed to fetch invalid projects:", err.Error())
        return
    }

    // delete all projects
    if count > 0 {
        err = query.Delete(new(ProjectMapping))
        if err != nil && err != storm.ErrNotFound {
            log.Errorln("failed to delete invalid projects:", err.Error())
            return
        }
        log.Infoln("successfully deleted", count, "project mappings")
    }

    // get all project mappings
    // in order to delete all builds which are not referenced anymore
    var mappings []ProjectMapping
    err = Get().All(&mappings)
    if err != nil && err != storm.ErrNotFound {
        log.Errorln("failed to fetch mapping:", err.Error())
        return
    }

    // construct the query
    unreferenced := make([]q.Matcher, len(mappings))
    for i, mapping := range mappings {
        unreferenced[i] = q.Not(q.Eq("Project", mapping.Id))
    }

    // find the number of builds to delete
    query = Get().Select(unreferenced...)
    count, err = query.Count(new(Build))
    if err != nil {
        log.Errorln("failed to fetch unreferenced builds:", err.Error())
        return
    }

    // delete all builds where no project mapping exists
    if count > 0 {
        err = query.Delete(new(Build))
        if err != nil && err != storm.ErrNotFound {
            log.Errorln("failed to delete unreferenced builds:", err.Error())
            return
        }
        log.Infoln("successfully deleted", count, "builds from the database")
    }

    log.Infoln("metadata purge finished successfully in", time.Since(start))
}
