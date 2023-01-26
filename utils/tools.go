package utils

import (
	"os/exec"
)

func ExecCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	out, err := command.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func WriteFile(filePath string, content []byte) error {
	cmd := exec.Command("echo", string(content), ">", filePath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(filePath string) (string, error) {
	cmd := exec.Command("cat", filePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
