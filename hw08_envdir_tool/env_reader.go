package main

import (
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.

type EnvValue struct {
	Value string

	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.

// Variables represented as files where filename is name of variable, file first line is a value.

func ReadDir(dir string) (Environment, error) {
	dirFs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	environment := make(map[string]EnvValue)

	for _, dirEntry := range dirFs {
		if !dirEntry.IsDir() && !strings.Contains(dirEntry.Name(), "=") {
			tmpBuff, err := os.ReadFile(dir + "/" + dirEntry.Name())
			if err != nil {
				return nil, err
			}

			if len(tmpBuff) == 0 {
				environment[dirEntry.Name()] = EnvValue{NeedRemove: true}
			} else {
				wholeString := string(tmpBuff)
				value := strings.ReplaceAll(strings.TrimRight(strings.Split(wholeString, "\n")[0], " "), "\x00", "\n")
				environment[dirEntry.Name()] = EnvValue{Value: value}
			}
		}
	}

	return environment, nil
}
