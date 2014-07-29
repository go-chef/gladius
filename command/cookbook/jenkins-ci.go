package cookbookcommand

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/bigkraig/go-gitlab/gitlab"
	"github.com/codegangsta/cli"
	"github.com/go-chef/gladius/app"
	"github.com/go-chef/gladius/lib"
)

type JenkinsCIContext struct {
	log *logrus.Logger
	cfg *app.Configuration

	ProjectID   int
	ProjectName string
	GroupName   string
	BranchName  string
}

const (
	versionCommitMessageRegex = `.*@(\d+)\.(\d+)\.(\d+).*`
	jenkinsCommitMessage      = "Tagged %d.%d.%d"
	jenkinsCommitMessageRegex = `Tagged \d+\.\d+\.\d+`
)

func JenkinsCICommand(env *app.Environment) cli.Command {
	c := &JenkinsCIContext{log: env.Log, cfg: env.Config}
	cmd := &cli.Command{
		Name:   "jenkins-ci",
		Usage:  "To be used through Jenkins only. Executes tests, tagging, and uploading to a Chef server.",
		Action: c.Run,
	}

	return *cmd
}

func parseJobName(jobName string) (projectName, groupName, branchName string, err error) {
	parts := strings.Split(jobName, "__")

	if len(parts) < 2 && len(parts) > 3 {
		err = errors.New(fmt.Sprintf("Unable to determine job name from %s", jobName))
		return
	}

	groupName = parts[0]
	projectName = parts[1]

	if len(parts) == 3 {
		branchName = parts[2]
	}

	return
}

func isJenkinsCommit(c *gitlab.ProjectCommit) bool {
	ok, _ := regexp.Match(jenkinsCommitMessageRegex, []byte(gitlab.Stringify(c.Title)))
	return ok
}

func isMergeCommit(c *gitlab.ProjectCommit) bool {
	return len(c.ParentIDs) > 1
}

// this could also detect if this is part of a merge request and fetch
// a version number from the MR comments, bump it in the metadata.rb and push it back up
func (j *JenkinsCIContext) Run(c *cli.Context) {
	// make this a little more generic so that we can also use github repos
	// and detect which type through the jenkins environment variables
	gitLabClient := lib.NewGitLabClient(j.cfg.APIURL, j.cfg.APISecret)
	log := j.log

	environment, err := lib.NewJenkinsEnvironment()
	if err != nil {
		log.Errorln(err)
		syscall.Exit(1)
	}
	env := environment

	j.ProjectName, j.GroupName, j.BranchName, err = parseJobName(env.JobName)
	if err != nil {
		log.Errorln(err)
		syscall.Exit(1)
	}

	cbMetadata, err := lib.NewCookbookFromMetadata(env.Workspace)
	if err != nil {
		log.Errorln(err)
		return
	}

	// Strip the kitchen platforms that are not supported by this cookbook
	platforms := []app.Platform{}
	for _, platform := range j.cfg.TestKitchenConfiguration.Platforms {
		for _, supported := range cbMetadata.Supports {
			if strings.HasPrefix(platform.Name, supported) {
				platforms = append(platforms, platform)
			}
		}
	}
	j.cfg.TestKitchenConfiguration.Platforms = platforms

	j.ProjectID, err = gitLabClient.FindProject(j.ProjectName, j.GroupName)
	if err != nil {
		log.Errorln(err)
		syscall.Exit(1)
	}

	commit, err := gitLabClient.FindCommit(env.GitCommit, j.ProjectID)
	if err != nil {
		log.Errorln(err)
		syscall.Exit(1)
	}
	log.Infoln(fmt.Sprintf("Testing commit %s (%s) for %s/%s",
		gitlab.Stringify(commit.Title), env.GitCommit[0:9], j.GroupName, j.ProjectName))

	if isJenkinsCommit(commit) {
		log.Infoln("Commit was done by Jenkins, skipping build")
		return
	}

	testSuite := &TestContext{
		cfg: j.cfg,
		log: log,
	}
	testSuite.Run(c)

	//  This is a regular commit and it has passed the tests, we can simply exit to Jenkins now
	if !isMergeCommit(commit) {
		log.Infoln("Passed tests.")
		return
	}

	// TODO: Locate merge request: impossible until merge request API exposes commit information
	// scan merge request commits for version change -- check merge request comments for version regex
	// if neither, then bump cookbook version and commit

	cookbook, err := lib.NewCookbook(gitLabClient, j.ProjectID, env.GitCommit)
	if err != nil {
		log.Errorln(err)
		syscall.Exit(1)
	}

	// TODO: Tag the release
	// Waiting on https://github.com/gitlabhq/gitlabhq/pull/7014

	// Validate that the cookbook version does not exist on the chef servers
	// may want to fail the build if the cookbook exists everywhere?
	for _, chefServer := range j.cfg.ChefServers {
		log.Infoln(fmt.Sprintf("Verifying that %s (%s) does not exist on %s", j.ProjectName,
			cookbook.Version(), chefServer.ServerURL))
		cb, err := chefServer.Client.Cookbooks.GetVersion(j.ProjectName, cookbook.Version())
		if cb.Name != "" {
			log.Warnln(fmt.Sprintf("Cookbook %s (%s) already exists on %s", j.ProjectName,
				cookbook.Version(), chefServer.ServerURL))
		}
		if err != nil {
			log.Errorln(err)
		}
	}

	// Berkshelf
	// The install and update commands dont get the custom berks configuration
	// This may be an issue if the cookbook expects to be able to read from our internal chef server
	// It's a complication when there are multiple chef servers to read cookbooks from
	if lib.NeedBerkshelfInstall(env.Workspace) {
		log.Infoln("Executing Berkshelf install step")
		if errs := lib.Execute("berks", "install"); errs > 0 {
			log.Errorln("Berkshelf install failed!")
			syscall.Exit(errs)
		}
	} else {
		log.Infoln("Executing Berkshelf update step")
		if errs := lib.Execute("berks", "update"); errs > 0 {
			log.Errorln("Berkshelf update failed!")
			syscall.Exit(errs)
		}
	}

	log.Infoln("Executing Berkshelf upload step")
	for _, chefServer := range j.cfg.ChefServers {
		file, err := lib.GenerateBerkshelfConfiguration(env.Workspace, &chefServer)
		if err != nil {
			log.Errorln(err)
			syscall.Exit(1)
		}
		log.Infoln("Uploading to", chefServer.ServerURL)
		if errs := lib.Execute("berks", "upload", "--config", file); errs > 0 {
			log.Errorln("Berkshelf upload failed!")
			syscall.Exit(errs)
		}
	}

	// Release the cookbook to the development environment
	log.Infoln(fmt.Sprintf("Releasing %s to 'development'", j.ProjectName))
	releaseCookbook := &ReleaseContext{
		cfg: j.cfg,
		log: log,
	}
	releaseCookbook.Do(j.ProjectName, cbMetadata.Version(), "development")
}
