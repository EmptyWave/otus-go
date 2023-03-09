package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	comm := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	comm.Stderr = os.Stderr
	comm.Stdin = os.Stdin
	comm.Stdout = os.Stdout

	comm.Env = prepareCmdEnv(env)

	err := comm.Run()
	if err != nil {
		returnCode = comm.ProcessState.ExitCode()
	}

	return
}

func prepareCmdEnv(env Environment) []string {
	for key, envVal := range env {
		if envVal.NeedRemove {
			os.Unsetenv(key)
			continue
		}

		os.Setenv(key, envVal.Value)
	}

	return os.Environ()
}
