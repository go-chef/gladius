package log

import (
	"github.com/Sirupsen/logrus"
)

// New creates a new logger with out default setup
func New() *logrus.Logger {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{}

	return log
}
