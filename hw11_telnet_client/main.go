package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := os.Args

	client := NewTelnetClient(args[1]+":"+args[2], *timeout, os.Stdin, os.Stdout)
	defer client.Close()
	err := client.Connect()
	if err != nil {
		fmt.Println("Could not connect to host")
		return
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		defer cancelCtx()
		err := client.Send()
		if err != nil {
			fmt.Println("Error while writing to host")
		}
	}()

	go func() {
		defer cancelCtx()
		err := client.Receive()
		if err != nil {
			fmt.Println("Error while reading from host")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	select {
	case <-sigChan:
	case <-ctx.Done():
		close(sigChan)
	}
}
