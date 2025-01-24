package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	timeOutFlag := flag.String("timeout", "0s", "connection timeout")
	flag.Parse()

	timeout, err := strconv.Atoi(strings.Split(*timeOutFlag, "s")[0])
	if err != nil {
		fmt.Println("Wrong timeout format")

		return
	}

	args := os.Args

	client := NewTelnetClient(args[1]+":"+args[2], time.Duration(timeout), os.Stdin, os.Stdout)

	err = client.Connect()
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
