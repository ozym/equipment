package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ozym/qdp"
)

func HostPort(ipaddr string, ipport string) (string, string) {
	s, p, err := net.SplitHostPort(ipaddr)
	if err == nil {
		return s, p

	}
	return ipaddr, ipport
}

func main() {

	// runtime settings
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "make noise")

	var serial bool
	flag.BoolVar(&serial, "serial", false, "recover instrument serial number details")

	var ipport string
	flag.StringVar(&ipport, "ipport", "5330", "Q330 port number to connect to")

	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 2*time.Second, "how long to wait")

	flag.Parse()

	results := make(map[string]interface{})

	for _, ipaddr := range flag.Args() {
		h, p := HostPort(ipaddr, ipport)
		if serial {
			s, err := qdp.ReadSerial(h, p, timeout)
			if err != nil {
				log.Fatal(err)
			}
			if s != nil {
				results[h] = s
			}
		} else {
			s, err := qdp.ReadSOH(h, p, timeout)
			if err != nil {
				log.Fatal(err)
			}
			if s != nil {
				results[h] = s
			}
		}

	}

	j, err := json.MarshalIndent(results, "", "  ")
	if err == nil {
		fmt.Println((string)(j))
	}
}
