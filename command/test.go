package command

import (
	"github.com/codegangsta/cli"
	"github.com/go-chef/gladius/app"
)

type testCommand struct {
	env *app.Environment
}

func TestCommand(env *app.Environment) cli.Command {
	c := &testCommand{env: env}
	cmd := &cli.Command{
		Name:   "test",
		Usage:  "Just an example of how to add a subcommand",
		Action: c.Run,
	}

	return *cmd
}

func (t *testCommand) Run(c *cli.Context) {
	println("This doesn't do anything")
}
