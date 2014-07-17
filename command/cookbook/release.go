package cookbookcommand

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/go-chef/gladius/app"
)

type releaseCookbookCommand struct {
	env *app.Environment
}

func ReleaseCommand(env *app.Environment) cli.Command {
	c := &releaseCookbookCommand{env: env}
	cmd := &cli.Command{
		Name:   "release",
		Usage:  "Releases a cookbook to a Chef environment",
		Action: c.Run,
	}

	return *cmd
}

/*
 * release <cookbook name> <cookbook version> <environment>
 *
 * verify cookbook version exists on all chef servers
 * if not, abort
 * fetch environment, update version, write environment
 * fetch all other environments
 *  - in any environment where cookbook is not pinned, pin it to 0.0.0
 */
func (t *releaseCookbookCommand) Run(c *cli.Context) {
	for _, chef := range t.env.Config.ChefServers {
		nodes, err := chef.Client.Nodes.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(nodes)
	}
}
