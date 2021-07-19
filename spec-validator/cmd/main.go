package main

import (
	"github.com/newrelic/infrastructure-agent/pkg/log"
	"golang.org/x/tools/go/ssa/interp/testdata/src/os"
)

const(
	specsPathEnvVar="SPECS"
)

func main(){
	log.Info("running spec-validator")
	paths:=os.Getenv(specsPathEnvVar)
	log.Info("validating specs in path %s",paths)
}