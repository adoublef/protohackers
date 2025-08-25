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
			Method string       `json:"method"`
			Number *json.Number `json:"number"`
		}

		if err := json.Unmarshal(sc.Bytes(), &p); err != nil {
			fmt.Fprintf(rwc, "MALFORMED\n")
			break
		}
		if p.Method != "isPrime" || p.Number == nil {
			// log.Printf("malformed payload %v", p)
			fmt.Fprintf(rwc, "MALFORMED\n")
			break
		}
		log.Printf("payload %v", p)

		var v struct {
			Method string `json:"method"`
			Prime  bool   `json:"prime"`
		}
		v.Method = p.Method

		n, err := p.Number.Int64()
		if err != nil {
			// _ = enc.Encode(v)
			fmt.Fprintf(rwc, "MALFORMED\n")
			break
		}

		if isPrime(int(n)) {
			v.Prime = true
		}
		_ = enc.Encode(v)
	}
	_ = sc.Err()
	// if err := sc.Err(); err != nil {
	// 	// fmt.Fprintf(rwc, "MALFORMED\n")
	// }
}

func isPrime(n int) bool {
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return n > 1
}
