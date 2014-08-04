package cookbookcommand

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/bigkraig/go-gitlab/gitlab"
	"github.com/codegangsta/cli"
	"github.com/go-chef/chef"
	"github.com/go-chef/gladius/app"
	"github.com/go-chef/gladius/lib"
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
	gitLabClient := lib.NewGitLabClient(r.cfg.APIURL, r.cfg.APISecret)

	// Release the cookbook to the autoReleaseEnvironment environment
	log.Infoln(fmt.Sprintf("Releasing %s to '%s'", cookbookName, environmentName))
	environmentsRepoID, err := gitLabClient.FindProject(gitlabEnvironmentProjectName, gitlabEnvironmentGroupName)
	if err != nil {
		log.Errorln(err)
	}

	gitEnvironments, _, err := gitLabClient.Projects.Tree(environmentsRepoID, "", "")
	if err != nil {
		log.Errorln(err)
		syscall.Exit(1)
	}
	for _, gitEnvironment := range *gitEnvironments {
		if *gitEnvironment.Type != "tree" {
			continue
		}

		gitFiles, _, err := gitLabClient.Projects.Tree(environmentsRepoID, *gitEnvironment.Name, "")
		if err != nil {
			log.Errorln(err)
			syscall.Exit(1)
		}

		for _, file := range *gitFiles {
			if !strings.HasSuffix(*file.Name, "json") {
				continue
			}

			sourceContents, _, err := gitLabClient.Projects.GetFileContents(environmentsRepoID, "master", *gitEnvironment.Name+"/"+*file.Name)
			if err != nil {
				log.Errorln(err)
				syscall.Exit(1)
			}

			env := &chef.Environment{}
			err = json.NewDecoder(sourceContents).Decode(&env)
			if err != nil {
				log.Errorln(fmt.Sprintf("Invalid json in %s: %s", *file.Name, err))
				syscall.Exit(1)
			}

			changed := false
			if env.Name == environmentName {
				if env.CookbookVersions[cookbookName] != cookbookVersion {
					changed = true
				}
				env.CookbookVersions[cookbookName] = cookbookVersion
			} else if env.CookbookVersions[cookbookName] == "" {
				env.CookbookVersions[cookbookName] = "0.0.0"
				changed = true
			}

			if !changed {
				continue
			}

			p := &gitlab.ProjectFileParameters{
				FilePath:      *gitEnvironment.Name + "/" + *file.Name,
				BranchName:    "master",
				CommitMessage: fmt.Sprintf("Released %s[%s] to %s", cookbookName, env.CookbookVersions[cookbookName], env.Name),
			}

			log.Infoln(fmt.Sprintf("Released %s[%s] to %s // %s", cookbookName, env.CookbookVersions[cookbookName], *gitEnvironment.Name, env.Name))

			content, err := json.MarshalIndent(&env, "", "  ")
			if err != nil {
				log.Errorln(err)
				syscall.Exit(1)
			}

			_, _, err = gitLabClient.Projects.UpdateFile(environmentsRepoID, *p, content)
			if err != nil {
				log.Errorln(err)
				syscall.Exit(1)
			}
		}
	}

}
