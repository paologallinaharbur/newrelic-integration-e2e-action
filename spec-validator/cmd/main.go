package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

const (
	specsPathEnvVar = "SPECS"
)

var log = logrus.New()

func main() {
	log.Info("running spec-validator")
	paths := os.Getenv(specsPathEnvVar)
	log.Infof("validating specs in path %s", paths)
}
