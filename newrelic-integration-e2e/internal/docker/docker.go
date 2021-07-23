package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	dockerComposeBin = "docker-compose"
)


func DockerComposeRun(path string,container string, envVars map[string]string) error {
	DockerComposeBuild(path,container)
	args := []string{"-f", path,"run"}
	for k, v := range envVars {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}
	args = append(args,  "-d", container)
	cmd := exec.Command(dockerComposeBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DockerComposeDown(path string) error {
	args := []string{"-f", path,"down","-v"}
	cmd := exec.Command(dockerComposeBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}


func DockerComposeBuild(path string,container string) error {
	args := []string{"-f", path,"build",  "--no-cache", container}
	fmt.Println(strings.Join(args," "))
	cmd := exec.Command(dockerComposeBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
