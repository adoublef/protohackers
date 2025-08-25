package main

import (
	"io"
	"net"
)

func main() {
	ln, _ := net.Listen("tcp", "0.0.0.0:8080")
	for {
		rwc, _ := ln.Accept()
		go io.Copy(rwc, rwc)
	}
}
