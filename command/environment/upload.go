package environmentcommand

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/go-chef/chef"
	"github.com/go-chef/gladius/app"
	"github.com/go-chef/gladius/lib"
)

type UploadContext struct {
	log *logrus.Logger
	cfg *app.Configuration
}

func UploadCommand(env *app.Environment) cli.Command {
	c := &UploadContext{log: env.Log, cfg: env.Config}
	cmd := &cli.Command{
		Name:        "upload",
		Description: "Uploads the environment(s) to the Chef server(s)",
		Usage:       "<environment.json file(s)>",
		Action:      c.Run,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "server, s",
				Usage: "target chef server"},
		},
	}

	return *cmd
}

/*
 * upload <environment.json>
 *
 */
func (r *UploadContext) Run(c *cli.Context) {
	if len(c.Args()) < 1 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return
	}

	filenames := c.Args()
	log := r.log
	errors := 0
	var GroupName string

	environment, err := lib.NewJenkinsEnvironment()
	if err == nil {
		_, GroupName, _, err = lib.ParseJenkinsJobName(environment.JobName)
		if err != nil {
			log.Errorln(err)
			syscall.Exit(1)
		}

		if r.cfg.GitLabConfiguration.ConfigurationGroup != GroupName {
			log.Infoln(fmt.Sprintf("Executed from a Jenkins build but not in the %s group.", r.cfg.GitLabConfiguration.ConfigurationGroup))
		}
	}

	for _, chefServer := range r.cfg.ChefServers {
		if c.String("server") != "" && !strings.Contains(chefServer.ServerURL, c.String("server")) {
			continue
		}
		for _, filename := range filenames {
			file, err := os.Open(filename)
			if err != nil {
				log.Errorln(fmt.Sprintf("Unable to open %s: %s", filename, err))
				errors += 1
				continue
			}
			defer file.Close()

			v := &chef.Environment{}
			err = json.NewDecoder(file).Decode(&v)
			if err != nil {
				log.Errorln(fmt.Sprintf("Invalid json in %s: %s", filename, err))
				errors += 1
				continue
			}

			if r.cfg.GitLabConfiguration.ConfigurationGroup != GroupName {
				continue
			}

			_, err = chefServer.Client.Environments.Get(v.Name)
			if err != nil {
				_, err = chefServer.Client.Environments.Create(v)
				if err != nil {
					log.Errorln(fmt.Sprintf("Error creating environment from %s: %s", filename, err))
					errors += 1
					continue
				}
				log.Infoln(fmt.Sprintf("Created the %s environment on %s", v.Name, chefServer.ServerURL))
			} else {
				_, err = chefServer.Client.Environments.Put(v)
				if err != nil {
					log.Errorln(fmt.Sprintf("Error updating environment from %s: %s", filename, err))
					errors += 1
					continue
				}
				log.Infoln(fmt.Sprintf("Updated the %s environment on %s", v.Name, chefServer.ServerURL))
			}
		}
	}
	if errors != 0 {
		os.Exit(errors)
	}
}
