package cookbookcommand

import (
	"fmt"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/go-chef/gladius/app"
)

type ReleaseContext struct {
	log *logrus.Logger
	cfg *app.Configuration
}

func ReleaseCommand(env *app.Environment) cli.Command {
	c := &ReleaseContext{log: env.Log, cfg: env.Config}
	cmd := &cli.Command{
		Name:        "release",
		Description: "Pins the cookbook version in the specified environment",
		Usage:       "<cookbook name> <cookbook version> <environment>",
		Action:      c.Run,
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
func (r *ReleaseContext) Run(c *cli.Context) {
	if len(c.Args()) < 3 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return
	}
	cookbookName := c.Args()[0]
	cookbookVersion := c.Args()[1]
	environmentName := c.Args()[2]

	r.Do(cookbookName, cookbookVersion, environmentName)
}

func (r *ReleaseContext) Do(cookbookName, cookbookVersion, environmentName string) {
	log := r.log
	for _, chefServer := range r.cfg.ChefServers {
		environments, err := chefServer.Client.Environments.List()
		if err != nil {
			log.Errorln(err)
			syscall.Exit(1)
		}

		for thisEnvironment, _ := range *environments {
			if thisEnvironment == "_default" {
				continue
			}

			chefEnvironment, err := chefServer.Client.Environments.Get(thisEnvironment)
			if err != nil {
				log.Errorln(err)
				syscall.Exit(1)
			}

			if thisEnvironment == environmentName {
				chefEnvironment.CookbookVersions[cookbookName] = cookbookVersion
			} else if chefEnvironment.CookbookVersions[cookbookName] == "" {
				chefEnvironment.CookbookVersions[cookbookName] = "0.0.0"
			} else {
				continue
			}

			log.Infoln(fmt.Sprintf("Pinning %s[%s] in %s on %s", cookbookName,
				chefEnvironment.CookbookVersions[cookbookName], thisEnvironment, chefServer.ServerURL))
			err = chefServer.Client.Environments.Put(chefEnvironment)
			if err != nil {
				log.Errorln("err", err)
				syscall.Exit(1)
			}
		}

	}
}
