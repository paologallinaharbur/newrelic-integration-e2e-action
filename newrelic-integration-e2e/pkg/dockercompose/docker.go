package dockercompose

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	dockerComposeBin = "docker-compose"
	dockerBin        = "docker"
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
	return cmd.Run()
}

func Down(path, containerName string) error {
	containerID, err := getContainerID(path, containerName)
	if err == nil {
		Logs(containerID)
	}
	args := []string{"-f", path, "down", "-v"}
	cmd := exec.Command(dockerComposeBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Build(path string, container string) error {
	args := []string{"-f", path, "build", "--no-cache", container}
	fmt.Println(strings.Join(args, " "))
	cmd := exec.Command(dockerComposeBin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Logs(containerID string) error {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("cntid: %s", containerID)
	args := []string{"logs", containerID}
	fmt.Println(strings.Join(args, " "))
	cmd := exec.Command(dockerBin, args...)
	stdout, err := cmd.Output()
	logrus.Debug("stdout")
	logrus.Debug(string(stdout))
	logrus.Debug(err)
	return err
}

func getContainerID(path, containerName string) (string, error) {
	args := []string{"-f", path, "ps", "-q", containerName}
	cmd := exec.Command(dockerComposeBin, args...)
	containerID, err := cmd.Output()
	return string(containerID), err
}
