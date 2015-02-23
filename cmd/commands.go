package cmd

import (
	"log"

	"github.com/go-chef/gladius/cmd/policy"
	"github.com/jawher/mow.cli"
)

func (c *CLI) RegisterCommands() {
	c.Command("policy", "Policy File related commands", func(cmd *cli.Cmd) {
		cmd.Command("install", "create a policyfile", pending)
		cmd.Command("update", "update a policyfile", policy.Update)
	})

	c.Command("upload", "Upload objects", func(cmd *cli.Cmd) {
		cmd.Command("cookbook", "Upload cookbook", pending)
		cmd.Command("role", "Upload Role", pending)
		cmd.Command("node", "Upload Node", pending)
		cmd.Command("environment", "Upload Environment", pending)
	})
}

func pending(cmd *cli.Cmd) {
	cmd.Action = func() {
		log.Fatal("Command not yet implemented!")
	}
}
