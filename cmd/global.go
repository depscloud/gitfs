package cmd

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	ConfigPath = "${HOME}/.gitfs/config.yml"
)

func fail(message string, data ...interface{}) {
	logrus.Errorf("[main] "+message, data...)
	os.Exit(1)
}
