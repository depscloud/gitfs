package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// ConfigPath defines where the default configuration file is stored.
	ConfigPath = "${HOME}/.gitfs/config.yml"
)

func fail(message string, data ...interface{}) {
	logrus.Errorf("[main] "+message, data...)
	os.Exit(1)
}
