package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/newrelic/newrelic-integration-e2e-action/spec-validator/pkg"
)

var log = logrus.New()

const (
	flagSpecPath = "spec_path"
	flagVerboseMode = "verbose_mode"
)

func main() {
	log.Info("running spec-validator")
	specsPathPtr := flag.String(flagSpecPath, "", "Relative path to the spec file")
	verboseModePtr := flag.Bool(flagVerboseMode, false, "If true the debug level is enabled")
	flag.Parse()
	specPath := *specsPathPtr
	if specPath == "" {
		os.Exit(1)
	}
	if *verboseModePtr{
		log.SetLevel(logrus.DebugLevel)
	}
	content,err:=ioutil.ReadFile(specPath)
	if err!=nil{
		log.Error(err)
		os.Exit(1)
	}
	log.Debug("parsing the content of the spec file")
	pkg.ParseSpecFile(content)
}

