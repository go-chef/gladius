package command

import (
	"github.com/codegangsta/cli"
	"github.com/go-chef/gladius/app"
	"github.com/go-chef/gladius/command/role"
)

func RoleCommand(env *app.Environment) cli.Command {
	cmd := &cli.Command{
		Name:  "role",
		Usage: "Role commands",
		Subcommands: []cli.Command{
			rolecommand.UploadCommand(env),
		},
	}

	return *cmd
}
