package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
)

func main() {
	ln, _ := net.Listen("tcp", "0.0.0.0:8080")
	serve(ln)
}

func serve(ln net.Listener) error {
	for {
		rwc, err := ln.Accept()
		if err != nil {
			return err
		}

		go handle(rwc)
	}
}

func handle(rwc net.Conn) {
	defer rwc.Close()

	var p struct {
		Method string
		Number *int
	}
	if err := json.NewDecoder(rwc).Decode(&p); err != nil {
		fmt.Fprintf(rwc, "MALFORMED\n")
		return
	}
	if p.Method != "isPrime" || p.Number == nil {
		fmt.Fprintf(rwc, "MALFORMED\n")
		return
	}

	var v struct {
		Method string `json:"method"`
		Prime  bool   `json:"prime"`
	}
	v.Method = p.Method
	e := json.NewEncoder(rwc)
	if isPrime(*p.Number) {
		v.Prime = true
	}
	_ = e.Encode(v)
}

func isPrime(n int) bool {
	for i := 2; i <= int((math.Sqrt(float64(n)))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return n > 1
}
