package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
			Number *float64 `json:"number,omitempty"` // ""
		}

		err := json.Unmarshal(sc.Bytes(), &p)
		if err != nil {
			fmt.Fprintf(rwc, "MALFORMED\n")
			break
		}
		if p.Method != "isPrime" || p.Number == nil {
			fmt.Fprintf(rwc, "MALFORMED\n")
			break
		}

		var v struct {
			Method string `json:"method"`
			Prime  bool   `json:"prime"`
		}
		v.Method = p.Method

		if isPrime(*p.Number) {
			v.Prime = true
		}
		_ = enc.Encode(v)
	}
	_ = sc.Err()
}

func isPrime(n float64) bool {
	if n == float64(int(n)) {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if int64(n)%int64(i) == 0 {
			return false
		}
	}
	return n > 1
}
