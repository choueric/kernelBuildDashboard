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
	"github.com/choueric/cmdmux"
)

var (
	buildProfile string
)

func initBuldCmd() {
	cmdmux.HandleFunc("/build", buildImageHandler)
	flagSet, err := cmdmux.FlagSet("/build")
	if err != nil {
		clog.Fatal(err)
	}
	flagSet.StringVar(&buildProfile, "p", "", "Specify profile by name or index.")

	cmdmux.HandleFunc("/build/image", buildImageHandler)
	cmdmux.SetFlagSet("/build/image", flagSet)

	cmdmux.HandleFunc("/build/modules", buildModulesHandler)
	cmdmux.SetFlagSet("/build/modules", flagSet)

	cmdmux.HandleFunc("/build/dtb", buildDtbHandler)
	cmdmux.SetFlagSet("/build/dtb", flagSet)
}

////////////////////////////////////////////////////////////////////////////////

func build_usage() {
	printTitle("- build [image|modules|dtb] [profile]")
	fmt.Printf("  Build various targets.")
	fmt.Printf(" Same as '$ kbdashboard make uImage' if target in config is uImage.\n")
}

////////////////////////////////////////////////////////////////////////////////

func image_usage() {
	printTitle("  - build image [profile]")
	fmt.Printf("    Build kernel images of [profile].\n")
	printDefOption("build")
}

func buildImageHandler(args []string, data interface{}) (int, error) {
	return wrap(build_image, args, data)
}

func build_image(args []string, config *Config) int {
	p, _ := getProfile(buildProfile, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}
	printCmd("build image", p.Name)
	return makeKernel(p, p.Target)
}

////////////////////////////////////////////////////////////////////////////////

func modules_usage() {
	printTitle("  - build modules [profile]")
	fmt.Printf("    Build and install modules of [profile].")
	fmt.Printf(" Same as '$ kbdashboard make modules' follwing\n")
	fmt.Printf("    '$ kbdashboard make modules_install'.\n")
}

func buildModulesHandler(args []string, data interface{}) (int, error) {
	return wrap(build_modules, args, data)
}

func build_modules(args []string, config *Config) int {
	p, _ := getProfile(buildProfile, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}

	printCmd("modules", p.Name)

	ret := makeKernel(p, "modules")
	if ret != 0 {
		clog.Fatalf("make modules failed.\n")
	}

	return makeKernel(p, "modules_install")
}

////////////////////////////////////////////////////////////////////////////////

func dtb_usage() {
	printTitle("  - build dtb [profile]")
	fmt.Printf("    Build dtb file specified in configration and install to 'BuildDir'.\n")
}

func build_dtb(args []string, config *Config) int {
	p, _ := getProfile(buildProfile, config)
	if p == nil {
		clog.Fatalf("can not find profile [%s]\n", args[0])
	}
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
	return wrap(build_dtb, args, data)
}
