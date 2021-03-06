package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
)

func doMakeKernelOpt(t string, n int, b, c, a, m, extra string) []string {
	cmdArgs := []string{}

	if t != "" {
		cmdArgs = append(cmdArgs, t)
	}

	if n > 0 {
		j := []string{"-j", strconv.Itoa(n)}
		cmdArgs = append(cmdArgs, strings.Join(j, ""))
	}

	if b != "" {
		output := []string{"O", b}
		cmdArgs = append(cmdArgs, strings.Join(output, "="))
	}

	if c != "" {
		cc := []string{"CROSS_COMPILE", c}
		cmdArgs = append(cmdArgs, strings.Join(cc, "="))
	}

	if a != "" {
		arch := []string{"ARCH", a}
		cmdArgs = append(cmdArgs, strings.Join(arch, "="))
	}

	if m != "" {
		installModPath := []string{"INSTALL_MOD_PATH", m}
		cmdArgs = append(cmdArgs, strings.Join(installModPath, "="))
	}

	if extra != "" {
		for _, p := range strings.Fields(extra) {
			cmdArgs = append(cmdArgs, p)
		}
	}

	return cmdArgs
}

func makeKernelOpt(p *Profile, target string) []string {
	return doMakeKernelOpt(target, p.ThreadNum, p.BuildDir, p.CrossComile,
		p.Arch, p.ModInstallDir, p.ExtraOpts)
}

func makeKernel(p *Profile, target string, w io.Writer, useMarker bool) error {
	if err := checkDirExist(p.BuildDir); err != nil {
		return err
	}
	if err := checkDirExist(p.ModInstallDir); err != nil {
		return err
	}

	cmdArgs := makeKernelOpt(p, target)
	logger.Println(cWrap(cGREEN, fmt.Sprintf("%v", cmdArgs)))

	cmd := exec.Command("make", cmdArgs...)
	if target == "kernelversion" {
		cmd.Dir = p.BuildDir
	} else {
		cmd.Dir = p.SrcDir
	}

	return pipeCmd(cmd, w, useMarker)
}

func kernelVersion(p *Profile) (string, error) {
	var buf bytes.Buffer

	err := makeKernel(p, "kernelversion", &buf, false)
	if err != nil {
		return "", err
	}
	version := buf.String()
	return version[0 : len(version)-1], nil
}

func localVersion(p *Profile) (string, error) {
	setlocalversion := path.Join(p.SrcDir, "scripts/setlocalversion")
	cmd := exec.Command(setlocalversion, p.SrcDir)
	cmd.Dir = p.BuildDir

	var buf bytes.Buffer
	err := pipeCmd(cmd, &buf, false)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func kernelFullVersion(p *Profile) (string, error) {
	version, err := kernelVersion(p)
	if err != nil {
		return "", err
	}

	local, err := localVersion(p)
	if err == nil {
		version += local
	}

	return version, nil
}

func configKernel(p *Profile, target string) error {
	if err := checkDirExist(p.BuildDir); err != nil {
		return err
	}
	if err := checkDirExist(p.ModInstallDir); err != nil {
		return err
	}

	cmdArgs := makeKernelOpt(p, target)
	args := []string{"make"}
	args = append(args, cmdArgs...)

	os.Chdir(p.SrcDir)
	return execCmd("make", args)
}

// execute command with Stdout and Stderr being piped in a new process.
// wait until this cmd finishes.
func pipeCmd(cmd *exec.Cmd, w io.Writer, useMarker bool) error {
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdoutMarker := cWrap(cGREEN, ">>")
	stderrMarker := cWrap(cRED, "!!")

	stdoutDone := make(chan struct{})
	scanner := bufio.NewScanner(stdoutReader)
	go func() {
		for scanner.Scan() {
			if useMarker {
				fmt.Fprintln(w, stdoutMarker, scanner.Text())
			} else {
				fmt.Fprintln(w, scanner.Text())
			}
		}
		stdoutDone <- struct{}{}
		logger.Printf("End of stdout goroutine. error: %v\n", scanner.Err())
	}()

	stderrDone := make(chan struct{})
	errScanner := bufio.NewScanner(stderrReader)
	go func() {
		for errScanner.Scan() {
			if useMarker {
				fmt.Fprintln(w, stderrMarker, errScanner.Text())
			} else {
				fmt.Fprintln(w, errScanner.Text())
			}
		}
		stderrDone <- struct{}{}
		logger.Printf("End of stderr goroutine. error: %v\n", errScanner.Err())
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}
	logger.Println("Start pipeCmd")

	<-stdoutDone
	<-stderrDone

	err = cmd.Wait()
	if err != nil {
		return err
	}

	logger.Println("End of pipeCmd")
	return nil
}

// execute command directly.
func execCmd(name string, argv []string) error {
	binary, err := exec.LookPath(name)
	if err != nil {
		return err
	}

	env := os.Environ()
	err = syscall.Exec(binary, argv, env)
	if err != nil {
		return err
	}

	return nil
}
