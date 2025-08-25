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

	for {
		var p struct {
			Method string       `json:"method"`
			Number *json.Number `json:"number"`
		}
		if err := json.NewDecoder(rwc).Decode(&p); err != nil {
			log.Printf("malformed payload %v", p)
			if err.Error() == "EOF" {
				break
			}
			fmt.Fprintf(rwc, "MALFORMED\n")
			return
		}
		log.Printf("input payload %v", p)
		if p.Method != "isPrime" || p.Number == nil {
			log.Printf("malformed payload %v", p)
			fmt.Fprintf(rwc, "MALFORMED\n")
			return
		}

		// {"method":"isPrime","prime":false}
		var v struct {
			Method string `json:"method"`
			Prime  bool   `json:"prime"`
		}
		v.Method = p.Method

		n, err := p.Number.Int64()
		log.Printf("input %d", n)
		if err != nil {
			p, _ := json.Marshal(v)
			log.Printf("float response %q", string(p))
			rwc.Write([]byte(string(p) + "\n"))
			return
		}

		if isPrime(int(n)) {
			v.Prime = true
		}
		b, _ := json.Marshal(v)
		log.Printf("ok response %q", string(b))
		rwc.Write([]byte(string(b) + "\n"))
	}
}

func isPrime(n int) bool {
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return n > 1
}
