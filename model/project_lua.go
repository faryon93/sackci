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
    "reflect"
    "errors"

    log "github.com/sirupsen/logrus"
    "github.com/yuin/gopher-lua"
    "github.com/faryon93/sackci/assets"
    "io/ioutil"
    "regexp"
    "github.com/faryon93/sackci/util"
)


// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    LUA_TAG = "lua"
)

var (
    // regular expressions
    ReFnAlias = regexp.MustCompile("^\\s*\\${(.*?)}\\s*$")
)


// ----------------------------------------------------------------------------------
//  public members
// ----------------------------------------------------------------------------------

// Creates a lua vm for this project.
func (p *Project) CreateLuaVm() (error) {
    // setup the lua vm
    p.lua = lua.NewState()

    // we want to apply the global lua contenxt
    // with predefined functions for the user
    err := loadGlobalLua(p.lua)
    if err != nil {
        log.Errorf("failed to load global lua context: %s", util.FirstLine(err.Error()))
    }

    // range over all members and find the lua enabled fields
    typ := reflect.TypeOf(*p)
    val := reflect.ValueOf(*p)
    for i := 0; i < typ.NumField(); i++ {
        // if the "lua" tag exists we sould process the script
        field := typ.Field(i)
        tagValue := field.Tag.Get(LUA_TAG)
        if len(tagValue) > 0 && field.Type.Kind() == reflect.String {
            // empty script can be omitted
            script := val.Field(i).String()
            if len(script) <= 0 {
                continue
            }

            // TODO: arguments more generic...
            // parse the configured script
            // e.g. apply fn aliases, ...
            src := ParseLuaScript(tagValue, "body, branch", script)

            // evaluate the lua script
            err := p.lua.DoString(src)
            if err != nil {
                log.Errorln("failed to execute lua script:", util.FirstLine(err.Error()))
                continue
            }

            // check if the specified function exits now
            fn := p.lua.GetGlobal(tagValue)
            if fn.Type() != lua.LTFunction {
                log.Errorf("function \"%s\" not found", tagValue)
                continue
            }
        }
    }

    return nil
}

// Executes the trigger filter script.
func (p *Project) EvalTriggerFilter(body string) (bool, error) {
    // if the function is not present in the global context
    // we mark all trigger events as okay
    fn := p.lua.GetGlobal("trigger_filter")
    if fn == nil || fn.Type() != lua.LTFunction {
        return true, nil
    }

    // push arguments on the stack
    p.lua.Push(fn.(*lua.LFunction))
    p.lua.Push(lua.LString(body))
    p.lua.Push(lua.LString(p.Branch))

    // call the lua function
    // TODO: implement errfunc
    err := p.lua.PCall(2, 1, nil)
    if err != nil {
        return false, err
    }

    // get the return value and make sure it has the right type
    top := p.lua.GetTop()
    returnValue :=  p.lua.Get(top)
    if returnValue.Type() != lua.LTBool {
        return false, errors.New("invalid return value")
    }

    // pop the return value from the stack
    // in order to prevent stack overflows
    p.lua.Pop(1)

    return lua.LVAsBool(returnValue), nil
}


// ----------------------------------------------------------------------------------
//  private functions
// ----------------------------------------------------------------------------------

// Applies the global lua file to the VM.
func loadGlobalLua(lua *lua.LState) (error) {
    file, err := assets.FS(false).Open("/lua/global.lua")
    if err != nil {
        return err
    }

    buf, err := ioutil.ReadAll(file)
    if err != nil {
        return err
    }

    return lua.DoString(string(buf))
}

// Parses the lua script in an configuration property.
// Applies all custom syntax.
func ParseLuaScript(fn string, args string, src string) (string) {
    // apply function aliases
    match := ReFnAlias.FindStringSubmatch(src)
    if len(match) == 2 {
        return "function " + fn + "(" + args + ") return " + match[1] +  "(" + args + ") end"
    }

    return src
}
