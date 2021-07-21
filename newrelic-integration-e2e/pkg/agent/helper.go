package agent

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

func removeDirectoryContent(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer func(){
		_ = d.Close()
	}()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = source.Close()
	}()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func(){
		_ = destination.Close()
	}()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func makeDirs(perm os.FileMode, dirs ...string) error {
	for i := range dirs {
		dir := dirs[i]
		if err := os.Mkdir(dir, perm); err != nil {
			return err
		}
	}
	return nil
}

func processTemplate(t *template.Template, vars map[string]interface{}, outputPath string) error {
	var templateOut bytes.Buffer
	if err := t.Execute(&templateOut, vars); err != nil {
		return err
	}
	content := templateOut.String()
	return ioutil.WriteFile(outputPath, []byte(content), 0777)
}