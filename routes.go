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
// imports
// --------------------------------------------------------------------------------------

import (
    "github.com/faryon93/sackci/rest"
)


// --------------------------------------------------------------------------------------
// global variables
// --------------------------------------------------------------------------------------

var routes = Routes{
    // REST Endpoints
    {"GET", "/api/v1/project", rest.ProjectList},
    {"GET", "/api/v1/project/{id}", rest.ProjectGet},
    {"GET", "/api/v1/project/{project}/history", rest.GetBuildHistory},
    {"GET", "/api/v1/project/{project}/env", rest.EnvGet},
    {"GET", "/api/v1/project/build/{id}", rest.GetBuild},
}
