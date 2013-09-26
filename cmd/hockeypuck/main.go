/*
   Hockeypuck - OpenPGP key server
   Copyright (C) 2012, 2013  Casey Marshall

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, version 3.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

// hockeypuck is an OpenPGP keyserver.
package main

import (
	"fmt"
	"launchpad.net/gnuflag"
	. "launchpad.net/hockeypuck"
	"os"
"path/filepath"
)

func die(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

type cmdHandler interface {
	Name() string
	Desc() string
	Flags() *gnuflag.FlagSet
	Main()
}

type subCmd struct {
	flags *gnuflag.FlagSet
}

func (c subCmd) Flags() *gnuflag.FlagSet { return c.flags }

var cmds []cmdHandler = []cmdHandler{
	newRunCmd(),
	newLoadCmd(),
	newDropConstraintsCmd(),
	newCreateConstraintsCmd(),
	newPbuildCmd(),
	newHelpCmd(),
	newVersionCmd()}

func main() {
	if len(os.Args) < 2 {
		newHelpCmd().Main()
		return
	}
	var cmdArgs []string
	if len(os.Args) > 2 {
		cmdArgs = os.Args[2:]
	}
	for _, cmd := range cmds {
		if cmd.Name() == os.Args[1] {
			if flags := cmd.Flags(); flags != nil {
				flags.Parse(false, cmdArgs)
			}
			cmd.Main()
			return
		}
	}
	newHelpCmd().Main()
}

type helpCmd struct {
	subCmd
}

func (c *helpCmd) Name() string { return "help" }

func (c *helpCmd) Desc() string { return "Display this help message" }

func (c *helpCmd) Main() {
	fmt.Printf(`Hockeypuck -- Public Keyserver
https://launchpad.net/hockeypuck

Hockeypuck is a public keyserver that supports the
HTTP Keyserver Protocol, as well as peering with SKS.

Basic commands:
`)
	for _, cmd := range cmds {
		fmt.Printf("  %s %s\t\t%s\n", filepath.Base(os.Args[0]), cmd.Name(), cmd.Desc())
	}
	os.Exit(1)
}

func newHelpCmd() *helpCmd {
	return new(helpCmd)
}

type versionCmd struct {
	subCmd
}

func (c *versionCmd) Name() string { return "version" }

func (c *versionCmd) Desc() string { return "Display Hockeypuck version information" }

func (c *versionCmd) Main() {
	fmt.Println(Version)
	os.Exit(0)
}

func newVersionCmd() *versionCmd {
	return new(versionCmd)
}

type configuredCmd struct {
	subCmd
	configPath string
}

func (c configuredCmd) Main() {
	if c.configPath != "" {
		if err := LoadConfigFile(c.configPath); err != nil {
			die(err)
		}
	} else {
		// Fall back on default empty config
		SetConfig("")
	}
}
