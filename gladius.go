package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var log = logrus.New()

func main() {
	app := cli.NewApp()
	app.Name = "gladius"
	app.Commands = GenerateCommands()
	app.Run(os.Args)
}
