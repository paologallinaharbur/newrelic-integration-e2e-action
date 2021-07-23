package docker_compose

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	dockerComposeBin = "docker-compose"
)

func Run(path string, container string, envVars map[string]string) error {
	if err := Build(path, container); err != nil {
		return err
	}
	args := []string{"-f", path, "run"}
	for k, v := range envVars {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}
	args = append(args, "-d", container)
	cmd := exec.Command(dockerComposeBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return  cmd.Run()
}

func Down(path string) error {
	Logs(path)
	args := []string{"-f", path, "down", "-v"}
	cmd := exec.Command(dockerComposeBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err:=cmd.Run()
	fmt.Println(cmd.Output())
	return  err
}

func Build(path string, container string) error {
	args := []string{"-f", path, "build", "--no-cache", container}
	fmt.Println(strings.Join(args, " "))
	cmd := exec.Command(dockerComposeBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err:=cmd.Run()
	fmt.Println(cmd.Output())
	return  err
}

func Logs(path string) error {
	args := []string{"-f", path, "logs"}
	fmt.Println(strings.Join(args, " "))
	cmd := exec.Command(dockerComposeBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err:=cmd.Run()
	fmt.Println(cmd.Output())
	return  err
}
