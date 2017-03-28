package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

const (
	cRED         = "\x1b[31;1m"
	cGREEN       = "\x1b[32;1m"
	cYELLOW      = "\x1b[33;1m"
	cEND         = "\x1b[0;m"
	level1Indent = "  "
	level2Indent = "      "
	level3Indent = "          "
)

func defMark() string {
	if gConfig.Color {
		return cRED + "*" + cEND
	} else {
		return "*"
	}
}

func cWrap(color string, str string) string {
	if gConfig == nil {
		return str
	}
	if gConfig.Color {
		return color + str + cEND
	} else {
		return str
	}
}

func cmdTitle(format string, def bool, v ...interface{}) {
	if def {
		fmt.Printf("%s %s\n", defMark(), cWrap(cGREEN, fmt.Sprintf(format, v...)))
	} else {
		fmt.Println(level1Indent + cWrap(cGREEN, fmt.Sprintf(format, v...)))
	}
}

func cmdInfo(format string, v ...interface{}) {
	fmt.Printf(level2Indent + fmt.Sprintf(format, v...))
}

func subcmdTitle(format string, def bool, v ...interface{}) {
	if def {
		fmt.Printf("    %s %s\n", defMark(), cWrap(cGREEN, fmt.Sprintf(format, v...)))
	} else {
		fmt.Println(level2Indent + cWrap(cGREEN, fmt.Sprintf(format, v...)))
	}
}

func subcmdInfo(format string, v ...interface{}) {
	fmt.Printf(level3Indent + fmt.Sprintf(format, v...))
}

func printCmd(cmd string, profile string) {
	logger.Printf("run '%s' for [%s]\n", cWrap(cGREEN, cmd),
		cWrap(cGREEN, profile))
}

// if OutputDir and ModInstallDir is relative, change it to absolute
// by adding SrcDir prefix.
func fixRelativeDir(p string, pre string) string {
	if p == "" {
		return ""
	}
	if !path.IsAbs(p) {
		p = path.Join(pre, p)
	}
	return p
}

func fixHomePath(p string) (ret string) {
	if p == "" {
		return ""
	}
	home := os.Getenv("HOME")

	if strings.HasPrefix(p, "$HOME") {
		ret = strings.Replace(p, "$HOME", home, 1)
	} else if strings.HasPrefix(p, "~") {
		ret = strings.Replace(p, "~", home, 1)
	} else {
		ret = p
	}
	return ret
}

func isNumber(str string) bool {
	if m, _ := regexp.MatchString("^[0-9]+$", str); !m {
		return false
	}
	return true
}

func checkFileExsit(p string) (bool, error) {
	_, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// check if exist, if not, create one
func checkDirExist(p string) error {
	if p == "" {
		return nil
	}

	exist, err := checkFileExsit(p)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	err = os.MkdirAll(p, os.ModeDir|0775)
	if err != nil {
		return err
	}
	return nil
}

func copyFileContents(src, dst string) error {
	fmt.Printf("copy '%s' -> '%s'\n", cWrap(cGREEN, src), cWrap(cGREEN, dst))
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	err = out.Sync()
	return nil
}
