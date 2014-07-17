package lib

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	metadataVersionRegex = `version.*'(\d+)\.(\d+)\.(\d+)'.*`
)

type Cookbook struct {
	majorVersion int
	minorVersion int
	patchVersion int
}

func (c *Cookbook) Version() string {
	return fmt.Sprintf("%d.%d.%d", c.majorVersion, c.minorVersion, c.patchVersion)
}

func (c *Cookbook) ParseVersionString(s string) (err error) {
	re := regexp.MustCompile(metadataVersionRegex)
	m := re.FindStringSubmatch(s)

	c.majorVersion, err = strconv.Atoi(m[1])
	if err != nil {
		return
	}

	c.minorVersion, err = strconv.Atoi(m[2])
	if err != nil {
		return
	}

	c.patchVersion, err = strconv.Atoi(m[3])
	if err != nil {
		return
	}

	return
}

func NewCookbook(c *GitLabClient, projectID int, gitCommit string) (*Cookbook, error) {
	cookbook := &Cookbook{}
	metadataContents, _, err := c.Projects.GetFileContents(projectID, gitCommit, "metadata.rb")
	if err != nil {
		return cookbook, err
	}
	scanner := bufio.NewScanner(metadataContents)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "version") {
			cookbook.ParseVersionString(line)
		}
	}
	return cookbook, nil
}
