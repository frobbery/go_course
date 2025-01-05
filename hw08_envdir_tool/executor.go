package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.

func RunCmd(cmd []string, env Environment) (returnCode int) {
	newEnv := makeNewEnv(env)
	runCmd := exec.Command(cmd[0], cmd[1:]...)
	runCmd.Env = newEnv
	runCmd.Stdin = os.Stdin
	runCmd.Stdout = os.Stdout

	log.Println(newEnv)

	log.Println("Command start err:", runCmd.Start())

	log.Println("Command wait err:", runCmd.Wait())

	return runCmd.ProcessState.ExitCode()
}

func makeNewEnv(env Environment) []string {

	newEnv := make([]string, 0)
	for _, val := range os.Environ() {
		keyVal := strings.Split(val, "=")
		if newValue, ok := env[keyVal[0]]; !ok {
			if !newValue.NeedRemove {
				log.Println(keyVal[0]+"="+newValue.Value)
				newEnv = append(newEnv, keyVal[0]+"="+newValue.Value)
			}
		} else {
			newEnv = append(newEnv, val)
		}
	}

	return newEnv
}
