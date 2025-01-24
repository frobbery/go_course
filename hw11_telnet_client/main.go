package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	timeout, err := strconv.Atoi(strings.Split(*flag.String("timeout", "0s", "connection timeout"), "s")[0])
	if err != nil {
		fmt.Println("Wrong timeout format")
		return
	}
	args := os.Args
	client := NewTelnetClient(args[1] + ":" + args[2], time.Duration(timeout), os.Stdin, os.Stdout);
	//os.Stdin.

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	for _ = range c {
		client.Close()
		return;
	}
}
