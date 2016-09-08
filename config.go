package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	ConfigDir = ".config/kbdashboard"
)

const DefaultConfig = `
{
	"item": [
	{
		"name":"demo",
		"thread_num":4,
		"output_dir":"./_build",
		"cross_compile":"arm-eabi-",
		"arch":"arm",
		"mod_install_dir":"./_build/mod",
		"src_dir":"/home/user/kernel"
	}
	]
}
`

type Item struct {
	Name          string `json:"name"`
	SrcDir        string `json:"src_dir"`
	Arch          string `json:"arch"`
	CrossComile   string `json:"cross_compile"`
	OutputDir     string `json:"output_dir"`
	ModInstallDir string `json:"mod_install_dir"`
	ThreadNum     int    `json:"thread_num"`
}

type Config struct {
	Items []*Item `json:"item"`
}

func (i *Item) String() string {
	return fmt.Sprintf("'%s'\nSrcDir\t: %s\nArch\t: %s\nCC\t: %s",
		i.Name, i.SrcDir, i.Arch, i.CrossComile)
}

func checkConfigDir(path string) {
	homeDir := os.Getenv("HOME")
	err := os.MkdirAll(homeDir+"/"+path, os.ModeDir|0777)
	if err != nil {
		log.Println("mkdir:", err)
	}
}

func checkConfigFile(path string) string {
	if path == "" {
		path = os.Getenv("HOME") + "/" + ConfigDir + "/config.json"
	}
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		log.Println("create an empty config file.")
		file, err := os.Create(path)
		_, err = file.Write([]byte(DefaultConfig))
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	} else if err != nil {
		log.Fatal(err)
	}

	return path
}

func ParseConfig(path string) (*Config, error) {
	checkConfigDir(ConfigDir)
	path = checkConfigFile(path)

	file, err := os.Open(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}