package configs

import (
	"fmt"
	"log"
	"os"
)

func MustGetEnvVar(envVarName string) string {
	val, err := GetEnvVar(envVarName)
	if err != nil {
		msg := fmt.Sprintf("mandatory environment variable %s not found", envVarName)
		log.Panicf(msg)
		os.Exit(1)
	}
	return val
}

// GetEnvVar retrieves the environment variable with the supplied name and fails
// if it is not able to do so.
func GetEnvVar(envVarName string) (string, error) {
	envVar := os.Getenv(envVarName)
	if envVar == "" {
		envErrMsg := fmt.Sprintf("the environment variable '%s' is not set", envVarName)
		return "", fmt.Errorf(envErrMsg)
	}
	return envVar, nil
}
