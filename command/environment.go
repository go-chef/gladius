package command

import (
	"github.com/codegangsta/cli"
	"github.com/go-chef/gladius/app"
	"github.com/go-chef/gladius/command/environment"
)

func EnvironmentCommand(env *app.Environment) cli.Command {
	cmd := &cli.Command{
		Name:  "environment",
		Usage: "Environment commands",
		Subcommands: []cli.Command{
			environmentcommand.UploadCommand(env),
		},
	}

	return *cmd
}
