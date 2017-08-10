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

// --------------------------------------------------------------------------------------
//  constants
// --------------------------------------------------------------------------------------

const (
    ENV_BUCKET = "env"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

func GetProjectEnv(project uint64) (map[string]string, error) {
    env := make(map[string]string, 0)
    return env, Get(ENV_BUCKET, project, &env)
}

func InsertProjectEnv(project uint64, key string, value string) (error) {
    env, _ := GetProjectEnv(project)
    env[key] = value

    return InsertId(ENV_BUCKET, project, env)
}