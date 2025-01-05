package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args

	env, err := ReadDir(args[1])
	if err != nil {
		log.Println("Error occurred while reading environment")

		return
	}

	returnCode := RunCmd(args[2:], env)

	log.Println("Return code ", returnCode)
}
