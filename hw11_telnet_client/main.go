package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := os.Args
	client := NewTelnetClient(args[2]+":"+args[3], *timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Println("Could not connect to host")
		return
	}
	defer client.Close()

	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		defer cancelCtx()
		err := client.Send()
		if err != nil {
			log.Println("Error while writing to host")
		}
	}()

	go func() {
		defer cancelCtx()
		err := client.Receive()
		if err != nil {
			log.Println("Error while reading from host")
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
