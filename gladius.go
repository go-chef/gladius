package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var log = logrus.New()

func main() {
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	app := cli.NewApp()
	app.Name = "gladius"
	app.Usage = "it does stuff"
	app.Commands = GenerateCommands()
	app.Run(os.Args)
}
