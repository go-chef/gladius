package lib

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Cookbook struct {
	majorVersion int
	minorVersion int
	patchVersion int
	Supports     []string
}

func (c *Cookbook) Version() string {
	return fmt.Sprintf("%d.%d.%d", c.majorVersion, c.minorVersion, c.patchVersion)
}

func (c *Cookbook) ParseVersionString(s string) (err error) {
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return errors.New(fmt.Sprintf("Unable to parse version from %s", s))
	}

	c.majorVersion, err = strconv.Atoi(parts[0])
	if err != nil {
		return
	}

	c.minorVersion, err = strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	c.patchVersion, err = strconv.Atoi(parts[2])
	if err != nil {
		return
	}

	return
}

func NewCookbookFromReader(r io.Reader) *Cookbook {
	cookbook := &Cookbook{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		spaceparts := strings.Split(line, " ")
		switch spaceparts[0] {
		case "version":
			parts := strings.Split(line, "'")
			if len(parts) > 1 {
				_ = cookbook.ParseVersionString(parts[1])
			}
		case "supports":
			parts := strings.Split(line, "'")
			if len(parts) > 1 {
				cookbook.Supports = append(cookbook.Supports, parts[1])
			}
		}
	}
	return cookbook
}

func NewCookbookFromMetadata(path string) (*Cookbook, error) {
	file, err := os.Open(fmt.Sprintf("%s%c%s", path, os.PathSeparator, "metadata.rb"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cookbook := NewCookbookFromReader(file)
	return cookbook, nil
}

func NewCookbook(c *GitLabClient, projectID int, gitCommit string) (*Cookbook, error) {
	metadataContents, _, err := c.Projects.GetFileContents(projectID, gitCommit, "metadata.rb")
	if err != nil {
		return nil, err
	}
	cookbook := NewCookbookFromReader(metadataContents)
	return cookbook, nil
}
