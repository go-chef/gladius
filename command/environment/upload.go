package environmentcommand

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/go-chef/chef"
	"github.com/go-chef/gladius/app"
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
	r.Do(c.Args())
}

func (r *UploadContext) Do(filenames []string) {
	log := r.log
	errors := 0
	for _, chefServer := range r.cfg.ChefServers {
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

			_, err = chefServer.Client.Environments.Get(v.Name)
			if err != nil {
				err = chefServer.Client.Environments.Create(v)
				if err != nil {
					log.Errorln(fmt.Sprintf("Error creating environment from %s: %s", filename, err))
					errors += 1
					continue
				}
				log.Infoln(fmt.Sprintf("Created the %s environment on %s", v.Name, chefServer.ServerURL))
			} else {
				err = chefServer.Client.Environments.Put(v)
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
