package main

import (
	"net"
	"testing"
	"time"

	"go.adoublef.dev/testing/is"
)

func Test_isPrime(t *testing.T) {
	// run server
	ln, err := net.Listen("tcp", "")
	is.OK(t, err) // net.Listen

	go serve(ln)
	t.Cleanup(func() { ln.Close() })

	conn, err := net.Dial("tcp", ln.Addr().String())
	is.OK(t, err) // net.Dial

	_, err = conn.Write([]byte(`{"method":"isPrime","number":-1}`))
	is.OK(t, err)

	// read response
	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second)) // avoid hanging forever
	n, err := conn.Read(buf)
	is.OK(t, err) // Read
	conn.Close()

	t.Logf("%s", string(buf[:n]))
}
