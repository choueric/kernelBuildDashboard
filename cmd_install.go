/*
 * Copyright (C) 2016 Eric Chou <zhssmail@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func installUsage() {
	printTitle("- install", false)
	fmt.Printf("  Execute the install script of current profile.\n")
}

// TODO: add arguments into the script.
func doInstall(args []string, config *Config) int {
	p, _ := getCurrentProfile(config)

	script := config.getInstallFilename(p)
	if checkFileExsit(script) == false {
		// create and edit script
		fmt.Printf("create install script: %s'%s'%s\n", CGREEN, script, CEND)
		file, err := os.OpenFile(script, os.O_RDWR|os.O_CREATE, 0775)
		checkError(err)
		str := fmt.Sprintf("#!/bin/sh\n\n# install script for profile '%s'", p.Name)
		_, err = file.Write([]byte(str))
		checkError(err)
		file.Close()
		return execCmd(config.Editor, []string{config.Editor, script})
	}

	printCmd("install", p.Name)
	fmt.Printf("    %s%s%s\n", CGREEN, script, CEND)
	// 1. cmd := exec.Command(script, "1", "2")
	cmd := exec.Command(script)
	// 2. cmd.Args = []string{script, "a", "b"}
	cmd.Dir = p.SrcDir
	return pipeCmd(cmd)
}

func installHandler(args []string, data interface{}) (int, error) {
	return wrap(doInstall, args, data)
}