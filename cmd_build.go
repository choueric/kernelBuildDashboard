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

	"github.com/choueric/clog"
)

func buildUsage() {
	cmdTitle("build [image|modules|dtb]", false)
	cmdInfo("Build various targets of kernel.\n")
	buildImageUsage()
	buildModulesUsage()
	buildDtbUsage()
	fmt.Printf("\n")
}

////////////////////////////////////////////////////////////////////////////////

func buildImageUsage() {
	subcmdTitle("build image", true)
	subcmdInfo("Build kernel images for current profile.\n")
	subcmdInfo("Equal to '$kbdashboard make uImage'.\n")
}

func doBuildImage(args []string, config *Config) int {
	p, _ := getCurrentProfile(config)
	printCmd("build image", p.Name)
	return makeKernel(p, p.Target)
}

func buildImageHandler(args []string, data interface{}) (int, error) {
	return wrap(doBuildImage, args, data)
}

////////////////////////////////////////////////////////////////////////////////

func buildModulesUsage() {
	subcmdTitle("build modules", false)
	subcmdInfo("Build and install modules for current profile.\n")
	subcmdInfo("Eqaul to '$ make modules' then '$ make modules_install'.\n")
}

func doBuildModules(args []string, config *Config) int {
	p, _ := getCurrentProfile(config)
	printCmd("modules", p.Name)

	ret := makeKernel(p, "modules")
	if ret != 0 {
		clog.Fatalf("make modules failed.\n")
	}

	return makeKernel(p, "modules_install")
}

func buildModulesHandler(args []string, data interface{}) (int, error) {
	return wrap(doBuildModules, args, data)
}

////////////////////////////////////////////////////////////////////////////////

func buildDtbUsage() {
	subcmdTitle("build dtb", false)
	subcmdInfo("Build 'dtb' file and install into 'BuildDir'.\n")
}

func doBuildDtb(args []string, config *Config) int {
	p, _ := getCurrentProfile(config)
	printCmd("build DTB", p.Name)

	if makeKernel(p, p.DTB) != 0 {
		clog.Fatalf("build DTB failed.\n")
	}

	src := p.OutputDir + "/arch/" + p.Arch + "/boot/dts/" + p.DTB
	dst := p.OutputDir + "/" + p.DTB

	if copyFileContents(src, dst) != nil {
		return 1
	} else {
		return 0
	}
}

func buildDtbHandler(args []string, data interface{}) (int, error) {
	return wrap(doBuildDtb, args, data)
}
