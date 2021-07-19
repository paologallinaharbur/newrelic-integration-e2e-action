package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

const (
	flagRootDir   = "root_dir"
	flagSpecsPath = "specs_path"
)

func main() {
	log.Info("running spec-validator")
	rootDirPtr := flag.String(flagRootDir, "", "Root  to the spec files")
	specsPathPtr := flag.String(flagSpecsPath, "", "Relative path to the spec files")
	flag.Parse()
	specsPath := *specsPathPtr
	rootDirPath := *rootDirPtr
	if specsPath == "" {
		os.Exit(1)
	}
	specFiles, err := findSpecFiles(rootDirPath, regexp.MustCompile(specsPath))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Infof("validating specs in path %s", specFiles)
}

func findSpecFiles(rootDir string, re *regexp.Regexp) ([]string, error) {
	files := []string{}
	filepath.Walk(rootDir, listFiles(re, files))
	log.Infof("found %[1]d spec files.\n", len(files))
	return files, nil
}

func listFiles(re *regexp.Regexp, files []string) func(fn string, fi os.FileInfo, err error) error {
	return func(fn string, fi os.FileInfo, err error) error {
		if re.MatchString(fn) == false {
			return errors.New("specs could not be found in given path")
		}
		if fi.IsDir() {
			dirFiles, err := ioutil.ReadDir(fi.Name())
			if err != nil {
				log.Fatal(err)
			}
			for _, f := range dirFiles {
				files = append(files, f.Name())
			}
			return nil
		}
		files = append(files, fn)
		return nil
	}
}
