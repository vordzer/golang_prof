package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.
	if len(cmd) == 0 {
		fmt.Println("haven't command")
		return -1
	}
	executor := exec.Command(cmd[0], cmd[1:]...)

	resEnv := make([]string, 0)
	for key, el := range env {
		resEnv = append(resEnv, key+"="+el.Value)
	}
	executor.Env = resEnv
	executor.Stdout = os.Stdout
	executor.Stdin = os.Stdin
	executor.Stderr = os.Stderr

	err := executor.Run()
	if exiterr, ok := err.(*exec.ExitError); ok {
		return exiterr.ExitCode()
	}
	return 0
}
