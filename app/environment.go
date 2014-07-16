package app

import "github.com/Sirupsen/logrus"

type Environment struct {
	Config *Configuration
	Log    *logrus.Logger
}
