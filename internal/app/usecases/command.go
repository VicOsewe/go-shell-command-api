package usecases

import "os/exec"

func ExecuteCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return string(exitError.Stderr), nil
		}
		return "", err
	}
	return string(output), nil
}
