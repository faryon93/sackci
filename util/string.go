package util
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
    "regexp"
)


// --------------------------------------------------------------------------------------
//  global variables
// --------------------------------------------------------------------------------------

var (
    reCredentialUrl = regexp.MustCompile("(.*)://(.*):(.*)@(.*)")
)


// --------------------------------------------------------------------------------------
//  public functions
// --------------------------------------------------------------------------------------

// Returns the short hash of a longer one.
func ShortHash(hash string) (string) {
    return hash[0:MinInt(12, len(hash))]
}

// Masks all creadentials in an url
func MaskCredentials(url string) (string) {
    match := reCredentialUrl.FindAllStringSubmatch(url, -1)
    if len(match) > 0 {
        return match[0][1] + "://****:****@" + match[0][4]
    }

    return url
}

// Checks if a letter is uppercase
func IsUpper(b byte) bool {
    return b < []byte{0x5a}[0]
}

// Returns true if all strings are empty.
func StrEmpty(strs ...string) bool {
    for _, str := range strs {
        if len(str) <= 0 {
            return true
        }
    }

    return false
}
