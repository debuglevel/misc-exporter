package main

import (
	"bytes"
	"log"
	"os/exec"
)

// ShellExecute executes a command using bash.
// Borrowed from https://stackoverflow.com/a/43246464/4764279
func ShellExecute(command string) (string, string, error) {
	log.Printf("Executing %v using shell...\n", command)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	log.Printf("Executed %v using shell\n", command)
	return stdout.String(), stderr.String(), err
}
