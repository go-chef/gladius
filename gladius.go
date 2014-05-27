package main

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/docopt/docopt-go"
	"github.com/kdar/factorlog"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// VERSION is the gladius version
const VERSION = "0.0.1"

// config type is the Native structure that holds our gladius config object
type config struct {
	ServerURL  string   `json:"server_url"`
	ClientName string   `json:"client_name"`
	KeyPath    string   `json:"client_key"`
	CookPaths  []string `json:"cook_paths"`
}

var Config = config{}

// setup our custom output logger
// TODO: not sure I like this logging lib, need to add logfile support possibly
var stderr = factorlog.New(os.Stderr, factorlog.NewStdFormatter(`%{Color "red" "ERROR"}%{Color "yellow" "WARN"}%{Color "green" "INFO"}%{Color "cyan" "DEBUG"}%{Color "blue" "TRACE"} %{SEVERITY}: %{Message}%{Color "reset"}`))

func main() {

	// Debatable weather this is worth doing or not.
	// TODO: benchmark/profile setting and leaving GOMAXPROCS default
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	args, _ := docopt.Parse(usage(), nil, true, VERSION, false, true)
	setlog(args["--loglevel"].(string))

	stderr.Debug("Args: ", spew.Sprint(args))
	configure(args["--config"].([]string))
	stderr.Debug("Config: ", spew.Sdump(Config))
}

// Configure finds, parses, and loads the config.json presented by the cli args. last file loaded wins.
func configure(files []string) {
	stderr.Debug("Config Paths: ", spew.Sprint(files))
	// Open and merge the configs
	for _, path := range files {
		path, err := filepath.Abs(path)
		if err != nil {
			stderr.Debugf("Couldn't get absolute path: %s\n\t%s", path, err.Error())
			continue
		}

		file, err := os.Open(path)
		if err != nil {
			stderr.Debugf("Couldn't open config: %s\n\t%s", path, err.Error())
			continue
		}
		stderr.Info("Loading Config: ", path)
		json.NewDecoder(file).Decode(&Config)
		if err != nil {
			stderr.Criticalf("Error processing file %s\n  %s", path, err.Error())
		}
		file.Close()
	}
}

// setlog Sets up our loglvel for the stderr logger
func setlog(level string) {
	switch strings.ToLower(level) {
	case "trace":
		stderr.SetMinMaxSeverity(factorlog.DEBUG, factorlog.ERROR)
	case "debug":
		stderr.SetMinMaxSeverity(factorlog.DEBUG, factorlog.ERROR)
	case "info":
		stderr.SetMinMaxSeverity(factorlog.INFO, factorlog.ERROR)
	case "warn":
		stderr.SetMinMaxSeverity(factorlog.WARN, factorlog.ERROR)
	case "error":
		stderr.SetSeverities(factorlog.ERROR)
	}
}

// usage  builds the usage text for docopt
func usage() string {
	// Find our User. So we can use it's home (in a platform independent way)
	usr, err := user.Current()
	if err != nil {
		stderr.Fatal(err)
	}

	usage := fmt.Sprintf(`
Usage: 
  gladius [-C <file>...][options] <action> <object> [NAME... | -] 

Options:
  -s <url>, --server <url>         Chef Server URL ex: https://myserver/orgname [default: https://localhost]
  -k <file>, --key <file>          Chef Client key file [default: /etc/chef/admin.pem]
  -c <name>, --client <name>       Chef client name [default: admin]
  -C <file>, --config <file>       Gladius config file to load can be specified multiple times last file wins. [default: %s/.gladius/config.json .gladius/config.json]
  -l <level>, --loglevel           Set output log levels: trace, debug, info, warn, error [default: info]
  --version   Output the version
  -h, --help  Get help text

Actions:
  download  Download an object frmo the server.
  show      Show object info.

Objects:
  cookbook, role, run_list, environment, data_bag 
`, usr.HomeDir)

	return usage
}
