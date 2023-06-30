package usecases

import (
	"fmt"
	"os/exec"
)

func ExecuteCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command does not exit with error: %v", err)
	}
	return string(output), nil
}
