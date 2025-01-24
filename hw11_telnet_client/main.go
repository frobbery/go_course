package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	//nolint:depguard
	"github.com/spf13/pflag"
)

func main() {
	var timeout time.Duration
	pflag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	pflag.Parse()

	args := os.Args

	client := NewTelnetClient(args[1]+":"+args[2], timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		fmt.Println("Could not connect to host")

		return
	}

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		readFromIn(client)
		wg.Done()
	}()

	go func() {
		writeToOut(client)
	}()

	go func() {
		listenForSignal()
	}()

	wg.Wait()
}

func readFromIn(client TelnetClient) {
	err := client.Send()
	if err != nil {
		fmt.Println("Error while reading from host")
	}
}

func writeToOut(client TelnetClient) {
	err := client.Receive()
	if err != nil {
		fmt.Println("Error while reading from host")
	}
}

func listenForSignal() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT)

	for range c {
		return
	}
}
