package main

import (
	"os"
	"runtime"

	"github.com/go-chef/gladius/cmd"
)

func main() {
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	app := cmd.NewCLI()
	app.RegisterCommands()
	app.Run(os.Args)
}
