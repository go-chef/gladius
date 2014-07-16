package lib

import (
	"errors"
	"reflect"
	"strconv"
	"syscall"
)

type JenkinsEnvironment struct {
	BuildCause              string `env:"BUILD_CAUSE"`
	BuildCauseManualTrigger bool   `env:"BUILD_CAUSE_MANUALTRIGGER"`
	BuildNumber             int    `env:"BUILD_NUMBER"`
	BuildTag                string `env:"BUILD_TAG"`
	BuildURL                string `env:"BUILD_URL"`
	ExecutorNumber          int    `env:"EXECUTOR_NUMBER"`
	GitBranch               string `env:"GIT_BRANCH"`
	GitCommit               string `env:"GIT_COMMIT"`
	GitURL                  string `env:"GIT_URL"`
	JenkinsHome             string `env:"JENKINS_HOME"`
	JenkinsURL              string `env:"JENKINS_URL"`
	JobName                 string `env:"JOB_NAME"`
	JobURL                  string `env:"JOB_URL"`
	Workspace               string `env:"WORKSPACE"`
}

func NewJenkinsEnvironment() (*JenkinsEnvironment, error) {
	environment := &JenkinsEnvironment{}
	c := reflect.ValueOf(environment)
	cType := reflect.TypeOf(*environment)

	for i := 0; i < c.Elem().NumField(); i++ {
		field := c.Elem().Field(i)
		value, _ := syscall.Getenv(cType.Field(i).Tag.Get("env"))
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(value, 10, 0)
			if err != nil {
				continue
			}
			field.SetInt(i)
		case reflect.String:
			field.SetString(value)
		case reflect.Bool:
			b, err := strconv.ParseBool(value)
			if err != nil {
				continue
			}
			field.SetBool(b)
		}
	}
	err := environment.validate()
	if err != nil {
		return environment, err
	}

	return environment, nil
}

func (e *JenkinsEnvironment) validate() error {
	switch {
	case e.GitCommit == "", e.GitBranch == "", e.JobName == "":
		return errors.New("JenkinsEnvironment is not configured, are you running from a Jenkins job?")
	}
	return nil
}
