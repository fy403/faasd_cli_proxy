package utils

import (
	"log"
	"os"
	"os/exec"
)

func ExecCommand(cmd string, args ...string) (string, error) {
	log.Printf("\t\tRuning cmd: %s args: %v", cmd, args)
	command := exec.Command(cmd, args...)
	out, err := command.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func WriteFile(filePath string, content []byte) error {
	log.Printf("\t\tWriting file: %s", filePath)
	err := os.WriteFile(filePath, content, 0666)
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(filePath string) ([]byte, error) {
	cmd := exec.Command("cat", filePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}
