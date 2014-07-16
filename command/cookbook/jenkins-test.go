package cookbookcommand

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"github.com/Sirupsen/logrus"
	"github.com/bigkraig/go-gitlab/gitlab"
	"github.com/codegangsta/cli"
	"github.com/go-chef/gladius/app"
	"github.com/go-chef/gladius/lib"
	"gopkg.in/yaml.v1"
)

type jenkinsTestContext struct {
	log *logrus.Logger
	cfg *app.Configuration
	lib.JenkinsEnvironment
	ProjectID   int
	ProjectName string
	GroupName   string

	client *gitlab.Client
}

func JenkinsTestCommand(env *app.Environment) cli.Command {
	c := &jenkinsTestContext{
		cfg: env.Config,
		log: env.Log,
	}
	cmd := &cli.Command{
		Name:   "test-suite",
		Usage:  "Runs a test suite against a cookbook",
		Action: c.Run,
	}

	return *cmd
}

const (
	versionCommitMessageRegex = `.*@(\d+)\.(\d+)\.(\d+).*`
	jenkinsCommitMessage      = "Tagged %d.%d.%d"
	jenkinsCommitMessageRegex = `Tagged \d+\.\d+\.\d+`
)

/*
 * jenkins-test takes no parameters and uses the jenkins environment variables to run
 *
 */
func (j *jenkinsTestContext) Run(c *cli.Context) {
	j.log.Infoln("jenkins test the cookbook: ", c.Args().First())

	j.client = gitlab.NewClient(j.cfg.APIURL, j.cfg.APISecret)

	environment, err := lib.NewJenkinsEnvironment()
	if err != nil {
		j.log.Errorln(err)
		return
	}
	j.JenkinsEnvironment = *environment

	err = j.parseJobName() // should check for an error here
	if err != nil {
		j.log.Errorln(err)
		syscall.Exit(1)
	}

	err = j.findProject()
	if err != nil {
		return
	}

	commit, err := j.findCommit()
	if err != nil {
		j.log.Errorln(err)
		syscall.Exit(1)
	}

	if isJenkinsCommit(commit) {
		j.log.Infoln("Commit was done by Jenkins, skipping build")
		return
	}

	quickErrors := 0

	// Run knife cookbook test

	// Run foodcritic
	j.log.Infoln("Executing Foodcritic")
	if errs := lib.Execute(j.Workspace, "foodcritic", ".", "-f", "any"); errs > 0 {
		j.log.Errorln("Foodcritic tests failed!")
		quickErrors += errs
	} else {
		j.log.Infoln("Foodcritic tests passed!")
	}

	// TODO: Chefspec ?

	// Run rubocop
	j.log.Infoln("Executing RuboCop")
	if errs := lib.Execute(j.Workspace, "rubocop"); errs > 0 {
		j.log.Errorln("RuboCop tests failed!")
		quickErrors += errs
	} else {
		j.log.Infoln("RuboCop tests passed!")
	}

	// If we failed any of the previous quick tests then lets abort before going through with the test kitchen
	if quickErrors > 0 {
		syscall.Exit(quickErrors)
	}

	// TODO: Run test kitchen
	err = j.generateTestKitchenConfiguration()
	j.log.Infoln("Executing Test Kitchen")
	if errs := lib.Execute(j.Workspace, "kitchen", "test"); errs > 0 {
		j.log.Errorln("Test Kitchen tests failed!")
		syscall.Exit(errs)
	} else {
		j.log.Infoln("Test Kitchen tests passed!")
	}
}

func (j *jenkinsTestContext) parseJobName() error {
	// TODO: Test if this worked
	f := strings.SplitN(j.JobName, "}-", 2)
	j.ProjectName = f[1]
	j.GroupName = f[0][1:]
	return nil
}

func (j *jenkinsTestContext) findProject() error {
	opts := &gitlab.SearchOptions{gitlab.ListOptions{Page: 1}}

	for {
		projects, resp, _ := j.client.Search.Projects(j.ProjectName, opts)
		for _, project := range *projects {
			if *project.Namespace.Name == j.GroupName && *project.Name == j.ProjectName {
				j.ProjectID = *project.ID
				return nil
			}
		}
		opts.Page = resp.NextPage
		if resp.NextPage == 0 {
			return errors.New(fmt.Sprintf("Unable to find %s in the %s group", j.ProjectName, j.GroupName))
		}
	}
}

func (j *jenkinsTestContext) findCommit() (*gitlab.ProjectCommit, error) {
	commit, _, err := j.client.Projects.GetCommit(j.ProjectID, j.GitCommit)
	if err != nil {
		j.log.Errorln(fmt.Sprintf("Error fetching commit %s in %s: %s", j.GitCommit, j.ProjectName, err))
		return commit, err
	}
	j.log.Infoln(fmt.Sprintf("Testing commit %s (%s) for %s/%s", gitlab.Stringify(commit.Title), j.GitCommit[0:9], j.GroupName, j.ProjectName))
	return commit, err
}

func isJenkinsCommit(c *gitlab.ProjectCommit) bool {
	ok, _ := regexp.Match(jenkinsCommitMessageRegex, []byte(gitlab.Stringify(c.Title)))
	return ok
}

func (j *jenkinsTestContext) generateTestKitchenConfiguration() error {
	kitchen := &app.TestKitchenConfiguration{
		Driver:      j.cfg.Driver,
		Provisioner: j.cfg.Provisioner,
		Platforms:   j.cfg.Platforms,
	}

	kitchenYAML, err := yaml.Marshal(&kitchen)
	if err != nil {
		j.log.Errorln(err)
		syscall.Exit(1)
	}

	filename := filepath.Join(j.Workspace, ".kitchen.local.yml")

	f, err := os.Create(filename)
	if err != nil {
		j.log.Errorln(err)
		syscall.Exit(1)
	}

	_, err = io.Copy(f, bytes.NewReader(kitchenYAML))
	if err != nil {
		j.log.Errorln(err)
		syscall.Exit(1)
	}
	f.Close()

	return nil
}
