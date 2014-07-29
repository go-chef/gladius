package cookbookcommand

import (
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/bigkraig/go-gitlab/gitlab"
	"github.com/codegangsta/cli"
	"github.com/go-chef/gladius/app"
	"github.com/go-chef/gladius/lib"
)

type TestContext struct {
	log *logrus.Logger
	cfg *app.Configuration

	client *gitlab.Client
}

func TestCommand(env *app.Environment) cli.Command {
	c := &TestContext{
		cfg: env.Config,
		log: env.Log,
	}
	cmd := &cli.Command{
		Name:   "test",
		Usage:  "Runs a test suite against a cookbook",
		Action: c.Run,
	}

	return *cmd
}

/*
 * test takes no parameters
 *
 */
func (t *TestContext) Run(c *cli.Context) {

	quickErrors := 0

	// Run foodcritic
	t.log.Infoln("Executing Foodcritic")
	if errs := lib.Execute("foodcritic", ".", "-f", "any"); errs > 0 {
		t.log.Errorln("Foodcritic tests failed!")
		quickErrors += errs
	} else {
		t.log.Infoln("Foodcritic tests passed!")
	}

	// TODO: Chefspec ?

	// Run rubocop
	t.log.Infoln("Executing RuboCop")
	if errs := lib.Execute("rubocop"); errs > 0 {
		t.log.Errorln("RuboCop tests failed!")
		quickErrors += errs
	} else {
		t.log.Infoln("RuboCop tests passed!")
	}

	// Run rspec
	t.log.Infoln("Executing RSpec")
	if errs := lib.Execute("rspec", "-c", "-fd"); errs > 0 {
		t.log.Errorln("RSpec tests failed!")
		quickErrors += errs
	} else {
		t.log.Infoln("RSpec tests passed!")
	}

	// If we failed any of the previous quick tests then lets abort before going through with the test kitchen
	if quickErrors > 0 {
		syscall.Exit(quickErrors)
	}

	// Run test kitchen
	err := lib.GenerateTestKitchenConfiguration(t.cfg)
	if err != nil {
		t.log.Errorln(err)
		syscall.Exit(1)
	}
	t.log.Infoln("Executing Test Kitchen")
	// make concurrency a config option?
	if errs := lib.Execute("kitchen", "test", "-c", "8"); errs > 0 {
		t.log.Errorln("Test Kitchen tests failed!")
		syscall.Exit(errs)
	} else {
		t.log.Infoln("Test Kitchen tests passed!")
	}
}
