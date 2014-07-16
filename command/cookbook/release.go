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
 * release-cookbook <cookbook name> <cookbook version> <environment>
 *
 * parse config file for chef servers, username, key
 * create clients
 * verify cookbook version exists on all chef servers
 * if not, abort
 * fetch push environment, update version, PUT envirnoment back
 * fetch all other environments
 *  - in any environment where cookbook is not pinned, pin it to 0.0.0
 */
func (t *releaseCookbookCommand) Run(c *cli.Context) {
	chefServers, err := t.env.Config.GenerateChefClients()
	if err != nil {
		t.env.Log.Errorln("Error generating chef client:", err)
		os.Exit(1)
	}

	for _, chef := range chefServers {
		nodes, err := chef.Nodes.List()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(nodes)
	}
}
