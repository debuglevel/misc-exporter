package main

import (
	"bytes"
	"os/exec"

	"go.uber.org/zap"
)

// ShellExecute executes a command using bash.
// Borrowed from https://stackoverflow.com/a/43246464/4764279
func ShellExecute(command string) (string, string, error) {
	zap.S().Debugf("Executing %v using shell...\n", command)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	zap.S().Debugf("Executed %v using shell\n", command)
	return stdout.String(), stderr.String(), err
}
