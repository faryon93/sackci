package ctx
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
//  types
// --------------------------------------------------------------------------------------

import (
    "github.com/faryon93/sackci/sse"
    "github.com/faryon93/sackci/config"
)


// --------------------------------------------------------------------------------------
//  global variables
// --------------------------------------------------------------------------------------

var (
    Feed *sse.Group
    Conf *config.Config
)


// --------------------------------------------------------------------------------------
//  initializer
// --------------------------------------------------------------------------------------

func Init() {
    // initialize the context variables
    Feed = sse.NewGroup("feed")
}


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

func RedirectToFeed(src chan sse.Event) {
    go func() {
        for v := range src {
            Feed.Publish(v)
        }
    }()
}
