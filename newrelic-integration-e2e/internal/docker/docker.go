package docker

import (
	"os"
	"os/exec"
)

const (
	dockerComposeBin = "docker-compose"
)

func DockerComposeUp(path string) error{
	cmd := exec.Command(dockerComposeBin, "-f", path, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}