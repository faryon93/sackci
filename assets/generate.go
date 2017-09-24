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

// Some generate rules
//go:generate lessc css/pipeline.less css/pipeline.css
//go:generate minify -a -r -o ../assets.min/ ./
//go:generate cp favicon.ico ../assets.min/
//go:generate esc -prefix=../assets.min -ignore (.*).go -o ../assets.go ../assets.min
//go:generate rm -rf ../assets.min
