package main

import (
	"encoding/json"
	"fmt"
	"log"
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
		Number *json.Number
	}
	if err := json.NewDecoder(rwc).Decode(&p); err != nil {
		log.Printf("malformed payload %v", p)
		fmt.Fprintf(rwc, "MALFORMED\n")
		return
	}
	log.Printf("input payload %v", p)
	if p.Method != "isPrime" || p.Number == nil {
		fmt.Fprintf(rwc, "MALFORMED\n")
		return
	}

	e := json.NewEncoder(rwc)

	var v struct {
		Method string `json:"method"`
		Prime  bool   `json:"prime"`
	}
	v.Method = p.Method

	n, err := p.Number.Int64()
	log.Printf("input %d", n)
	if err != nil {
		log.Printf("float response %v", v)
		_ = e.Encode(v)
		return
	}

	if isPrime(int(n)) {
		v.Prime = true
	}
	log.Printf("ok response %v", v)
	_ = e.Encode(v)
}

func isPrime(n int) bool {
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return n > 1
}
