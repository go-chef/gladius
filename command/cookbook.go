package command

import (
	"github.com/codegangsta/cli"
	"github.com/go-chef/gladius/app"
	"github.com/go-chef/gladius/command/cookbook"
)

func CookbookCommand(env *app.Environment) cli.Command {
	cmd := &cli.Command{
		Name:  "cookbook",
		Usage: "Cookbook commands",
		Subcommands: []cli.Command{
			cookbookcommand.TestCommand(env),
			cookbookcommand.ReleaseCommand(env),
			cookbookcommand.JenkinsCICommand(env),
		},
	}

	return *cmd
}
