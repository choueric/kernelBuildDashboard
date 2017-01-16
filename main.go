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

	"github.com/choueric/clog"
	"github.com/choueric/cmdmux"
)

func usageHandler(args []string, data interface{}) (int, error) {
	fmt.Printf("Usage:\n")
	usageList()
	return 1, nil
}

func main() {
	clog.SetFlags(clog.Lshortfile | clog.LstdFlags | clog.Lcolor)

	if len(os.Args) >= 2 && os.Args[1] == "dump" {
		getConfig(true)
		return
	}

	config := getConfig(false)

	// TODO: use wrap
	cmdmux.HandleFunc("/", usageHandler)

	initListCmd()
	initChooseCmd()
	initEditCmd()
	initConfigCmd()
	initBuldCmd()
	initInstallCmd()
	initMakeCmd()

	cmdmux.HandleFunc("/help", usageHandler)

	ret, err := cmdmux.Execute(config)
	if err != nil {
		clog.Warn("Execute error:", err)
		os.Exit(0)
	}
	os.Exit(ret)
}
