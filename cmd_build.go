package main

import (
	"fmt"
	"io"
)

func buildUsage(w io.Writer, m *helpMap) {
	defaultHelp(w, m)
	fmt.Printf("\n")
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

func buildImageHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}
	printCmd("build image", p.Name)
	return 0, makeKernel(p, p.Target)
}

////////////////////////////////////////////////////////////////////////////////

func buildModulesUsage() {
	subcmdTitle("build modules", false)
	subcmdInfo("Build and install modules for current profile.\n")
	subcmdInfo("Eqaul to '$ make modules' then '$ make modules_install'.\n")
}

func buildModulesHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}
	printCmd("build modules", p.Name)

	err = makeKernel(p, "modules")
	if err != nil {
		return 0, err
	}

	return 0, makeKernel(p, "modules_install")
}

////////////////////////////////////////////////////////////////////////////////

func buildDtbUsage() {
	subcmdTitle("build dtb", false)
	subcmdInfo("Build 'dtb' file and install into 'BuildDir'.\n")
}

func buildDtbHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}
	printCmd("build DTB", p.Name)

	if err := makeKernel(p, p.DTB); err != nil {
		return 0, err
	}

	src := p.BuildDir + "/arch/" + p.Arch + "/boot/dts/" + p.DTB
	dst := p.BuildDir + "/" + p.DTB

	if err := copyFileContents(src, dst); err != nil {
		return 0, err
	}
	return 0, nil
}
