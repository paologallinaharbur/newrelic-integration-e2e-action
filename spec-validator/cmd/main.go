package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/newrelic/newrelic-integration-e2e-action/spec-validator/pkg"
)

var log = logrus.New()

const (
	flagSpecPath = "spec_path"
)

func main() {
	log.Info("running spec-validator")
	specsPathPtr := flag.String(flagSpecPath, "", "Relative path to the spec file")
	flag.Parse()
	specPath := *specsPathPtr
	if specPath == "" {
		os.Exit(1)
	}
	parent:=filepath.Base(specPath)
	fmt.Println(parent)
	var files []string
	err := filepath.Walk(parent, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
	content,err:=ioutil.ReadFile(specPath)
	if err!=nil{
		log.Error(err)
		os.Exit(1)
	}
	pkg.ParseSpecFile(content)
}

