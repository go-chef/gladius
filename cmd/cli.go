package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-chef/gladius/log"
	"github.com/jawher/mow.cli"
)

// VERSION is the gladius version
const VERSION = "0.0.1"

var loger = log.New()

type CLI struct {
	*cli.Cli
	*logrus.Logger
}

// config type is the Native structure that holds our gladius config object
type config struct {
	ServerURL  string   `json:"server_url"`
	ClientName string   `json:"client_name"`
	KeyPath    string   `json:"client_key"`
	CookPaths  []string `json:"cook_paths"`
}

var Config = config{}

func NewCLI() (c *CLI) {
	c = &CLI{cli.App("gladius", "Golang chef cli"), loger}

	// TODO: make this windows friendly
	defaultConfigDirs := []string{
		"/etc/gladius/config.json",
		"~/.gladius/config.json",
		".gladius/config.json",
	}

	confFiles := c.Strings(cli.StringsOpt{
		Name:   "c config",
		Value:  defaultConfigDirs,
		EnvVar: "GLADIUS_CONFIG",
		Desc:   "Locations to scan for config file",
	})

	c.BoolOpt("d debug", false, "enable debug output")
	c.BoolOpt("v verbose", false, "enable verbose output")
	configure(confFiles)

}

// Configure finds, parses, and loads the config.json presented by the cli args. last file loaded wins.
func configure(files *[]string) {
	loger.Debug("Config Paths: ", spew.Sprint(files))
	// Open and merge the configs
	for _, path := range *files {
		path, err := filepath.Abs(path)
		if err != nil {
			loger.Debugf("Couldn't get absolute path: %s\n\t%s", path, err.Error())
			continue
		}

		file, err := os.Open(path)
		if err != nil {
			loger.Debugf("Couldn't open config: %s\n\t%s", path, err.Error())
			continue
		}
		defer file.Close()

		loger.Info("Loading Config: ", path)
		err = json.NewDecoder(file).Decode(&Config)
		if err != nil {
			loger.Errorf("Error processing file %s\n  %s", path, err)
			os.Exit(1)
		}
	}
}
