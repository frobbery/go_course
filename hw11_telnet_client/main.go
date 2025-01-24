package main

import (
	"bufio"
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

		wg.Done()
	}()

	go func() {
		listenForSignal()

		wg.Done()
	}()

	wg.Wait()
}

func readFromIn(client TelnetClient) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		err := client.Send()
		if err != nil {
			fmt.Println("Error while reading from host")
			break
		}

		if !scanner.Scan() {
			fmt.Println("EOF in stdin")

			break
		}
	}
}

func writeToOut(client TelnetClient) {
	for {
		err := client.Receive()
		if err != nil {
			fmt.Println("Error while reading from host")

			break
		}
	}
}

func listenForSignal() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT)

	for range c {
		return
	}
}
