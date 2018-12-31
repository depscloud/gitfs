package cmd

import (
	"github.com/sirupsen/logrus"
	"os"
)

func fail(message string, data ...interface{}) {
	logrus.Errorf("[main] "+message, data...)
	os.Exit(1)
}
