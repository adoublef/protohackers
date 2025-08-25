package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
)

func main() {
	ln, _ := net.Listen("tcp", "0.0.0.0:8080")
	defer ln.Close()
	serve(ln)
}

func serve(ln net.Listener) error {
	for {
		rwc, err := ln.Accept()
		if err != nil {
			continue
		}

		go handle(rwc)
	}
}

func handle(rwc net.Conn) {
	defer rwc.Close()

	enc := json.NewEncoder(rwc)

	sc := bufio.NewScanner(rwc)
	for sc.Scan() {
		var p struct {
			Method string   `json:"method"`
			Number *float64 `json:"number,omitempty"`
		}

		toString := func() string {
			if p.Number != nil {
				return fmt.Sprintf("method=%s,number=%s", p.Method, (p.Number))
			} else {
				return fmt.Sprintf("method=%s", p.Method)
			}
		}

		err := json.Unmarshal(sc.Bytes(), &p)
		if err != nil {
			log.Printf("mlfrm -- %v", toString())
			fmt.Fprintf(rwc, "MALFORMED\n")
			break
		}
		if p.Method != "isPrime" || p.Number == nil {
			log.Printf("inval -- %v", toString())
			fmt.Fprintf(rwc, "MALFORMED\n")
			break
		}
		log.Printf("input -- %v", toString())

		var v struct {
			Method string `json:"method"`
			Prime  bool   `json:"prime"`
		}
		v.Method = p.Method

		if n := *p.Number; !(n == float64(int(n))) {
			_ = enc.Encode(v)
			continue
		}

		if isPrime(int64(*p.Number)) {
			v.Prime = true
		}
		_ = enc.Encode(v)
	}
	_ = sc.Err()
}

func isPrime(n int64) bool {
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%int64(i) == 0 {
			return false
		}
	}
	return n > 1
}
