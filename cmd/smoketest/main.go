package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:8080")
	fatal(err)

	fatal(serve(ln))
}

func serve(ln net.Listener) error {
	defer ln.Close()
	for {
		rwc, err := ln.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept: %w", err)
		}

		go handle(rwc)
	}
}

func handle(rwc net.Conn) {
	defer rwc.Close()

	_, _ = io.Copy(rwc, rwc)
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
