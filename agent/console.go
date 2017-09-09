package agent
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
    "strings"
)


// ----------------------------------------------------------------------------------
//  types
// ----------------------------------------------------------------------------------

type ConsoleOutput struct {
    // public members
    Callback func(string)

    // private members
    buffer string
}


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

func (w *ConsoleOutput) Write(p []byte) (int, error) {
    w.buffer += string(p)
    lines := strings.Split(w.buffer, "\n")
    last := len(lines) - 1

    if last >= 0 {
        for i := 0; i < last; i++ {
            w.Callback(clean(lines[i]))
        }

        if strings.HasSuffix(lines[last], "\n") {
            w.Callback(clean(lines[last]))
            w.buffer = ""
        } else {
            w.buffer = lines[last]
        }
    }

    return len(p), nil
}


// ----------------------------------------------------------------------------------
//  private functions
// ----------------------------------------------------------------------------------

func clean(str string) (string) {
    // TODO: better handling of ansi terminal codes
    str = strings.Replace(str, "\u001b[K", "", -1)

    parts := strings.Split(str, "\r")
    if len(parts) >= 2 {
        return parts[len(parts) - 2]
    } else {
        return str
    }
}