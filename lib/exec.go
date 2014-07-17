package lib

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

func Execute(command ...string) int {
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Infoln("...running", strings.Join(command, " "))

	var cmd *exec.Cmd

	if len(command) > 1 {
		cmd = exec.Command(command[0], command[1:]...)
	} else {
		cmd = exec.Command(command[0])
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		log.Errorln(err)
		return 1
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			log.Infoln(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Errorln(scanner.Text())
		}
	}()

	err = cmd.Wait()
	if err != nil {
		log.Errorln(err)
		return 1
	}

	return 0
}
