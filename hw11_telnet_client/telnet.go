package main

import (
	"errors"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	con     net.Conn
}

func (c *client) Connect() error {
	con, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.con = con
	return nil
}

func (c *client) Receive() error {
	tmp := make([]byte, 256)
	for {
		n, err := c.con.Read(tmp)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			}
			break
		}
		_, err = c.out.Write(tmp[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *client) Send() error {
	tmp := make([]byte, 256)
	for {
		n, err := c.in.Read(tmp)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("Stop reading")
				return nil
			}
			return err
		}
		_, err = c.con.Write(tmp[:n])
		if err != nil {
			return err
		}
	}
}

func (c *client) Close() error {
	return c.con.Close()
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
