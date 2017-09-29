package assets
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
//  generate rules
// ----------------------------------------------------------------------------------
//go:generate lessc css/pipeline.less css/pipeline.css
//go:generate minify -a -r -o ../assets.min/ ./
//go:generate cp favicon.ico ../assets.min/
//go:generate esc -pkg assets -prefix=../assets.min -ignore (.*).go -o fs.go ../assets.min
//go:generate rm -rf ../assets.min


// ----------------------------------------------------------------------------------
//  public functions
// ----------------------------------------------------------------------------------

// Returns true if file exists in the embedded asset
// filesystem used to serve static http content.
func FileExists(file string) bool {
    entry, exists := _escData[file]
    if !exists {
        return false
    }

    // we only want files to be known as existing
    return !entry.isDir
}
