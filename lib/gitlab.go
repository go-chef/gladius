package lib

import (
	"errors"
	"fmt"

	"github.com/bigkraig/go-gitlab/gitlab"
)

type GitLabClient struct {
	*gitlab.Client
}

func NewGitLabClient(APIUrl, APISecret string) *GitLabClient {
	client := gitlab.NewClient(APIUrl, APISecret)
	return &GitLabClient{client}
}

func (g *GitLabClient) FindProject(projectName, groupName string) (projectID int, err error) {
	opts := &gitlab.SearchOptions{gitlab.ListOptions{Page: 1}}

	for {
		projects, resp, _ := g.Search.Projects(projectName, opts)
		for _, project := range *projects {
			if *project.Namespace.Name == groupName && *project.Name == projectName {
				projectID = *project.ID
				return
			}
		}
		opts.Page = resp.NextPage
		if resp.NextPage == 0 {
			return 0, errors.New(fmt.Sprintf("Unable to find %s in the %s group", projectName, groupName))
		}
	}
}

func (g *GitLabClient) FindCommit(gitCommit string, projectID int) (*gitlab.ProjectCommit, error) {
	commit, _, err := g.Projects.GetCommit(projectID, gitCommit)
	if err != nil {
		return commit, err
	}
	return commit, err
}
