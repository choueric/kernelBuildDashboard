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
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/choueric/clog"
)

func createDir(p string) {
	err := os.MkdirAll(p, os.ModeDir|0777)
	if err != nil {
		clog.Printf("mkdir %s failed: %v\n", p, err)
	}
}

func makeKernelOpt(p *Profile, target string) []string {
	cmdArgs := []string{}

	if target != "" {
		cmdArgs = append(cmdArgs, target)
	}

	if p.ThreadNum > 0 {
		j := []string{"-j", strconv.Itoa(p.ThreadNum)}
		cmdArgs = append(cmdArgs, strings.Join(j, ""))
	}

	if p.OutputDir != "" {
		output := []string{"O", p.OutputDir}
		cmdArgs = append(cmdArgs, strings.Join(output, "="))
		createDir(p.OutputDir)
	}

	if p.CrossComile != "" {
		cc := []string{"CROSS_COMPILE", p.CrossComile}
		cmdArgs = append(cmdArgs, strings.Join(cc, "="))
	}

	if p.Arch != "" {
		arch := []string{"ARCH", p.Arch}
		cmdArgs = append(cmdArgs, strings.Join(arch, "="))
	}

	if p.ModInstallDir != "" {
		installModPath := []string{"INSTALL_MOD_PATH", p.ModInstallDir}
		cmdArgs = append(cmdArgs, strings.Join(installModPath, "="))
		createDir(p.ModInstallDir)
	}

	return cmdArgs
}

func makeKernel(p *Profile, target string) int {
	cmdArgs := makeKernelOpt(p, target)

	cmd := exec.Command("make", cmdArgs...)
	cmd.Dir = p.SrcDir

	fmt.Printf("    %s%v%s\n", CGREEN, cmdArgs, CEND)

	return runCmd(cmd)
}

func configKernel(p *Profile, target string) int {
	cmdArgs := makeKernelOpt(p, target)
	args := []string{"make"}
	args = append(args, cmdArgs...)

	os.Chdir(p.SrcDir)
	return execCmd("make", args)
}

func runCmd(cmd *exec.Cmd) int {
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		clog.Println("Error creating StdoutPipe for Cmd:", err)
		return 1
	}

	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		clog.Println("create stderrPipe:", err)
		return 2
	}

	scanner := bufio.NewScanner(stdoutReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s>>%s %s\n", CGREEN, CEND, scanner.Text())
		}
	}()

	errScanner := bufio.NewScanner(stderrReader)
	go func() {
		for errScanner.Scan() {
			fmt.Printf("%s!!%s %s\n", CRED, CEND, errScanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		clog.Println("Error starting Cmd:", err)
		return 3
	}

	err = cmd.Wait()
	if err != nil {
		clog.Println("Error waiting for Cmd:", err)
		return 4
	}

	return 0
}

func execCmd(name string, argv []string) int {
	binary, err := exec.LookPath(name)
	if err != nil {
		clog.Fatal(err)
	}

	env := os.Environ()

	err = syscall.Exec(binary, argv, env)
	if err != nil {
		clog.Fatal(err)
	}

	return 0
}
