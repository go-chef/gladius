package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

import (
	"github.com/go-chef/gladius/app"
	"github.com/go-chef/gladius/command"
)

func GenerateCommands() []cli.Command {
	var Commands []cli.Command

	cfg, err := app.ReadConfiguration()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	env := &app.Environment{
		Config: cfg,
		Log:    log,
	}

	Commands = []cli.Command{
		command.CookbookCommand(env),
		command.TestCommand(env),
	}

	return Commands
}
