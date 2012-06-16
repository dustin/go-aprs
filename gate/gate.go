package main

import (
	"flag"
	"fmt"
	"log"
	"net/textproto"
	"os"

	"github.com/dustin/aprs.go"
)

var call, pass, filter, server string

func init() {
	flag.StringVar(&server, "server", "second.aprs.net:14580", "APRS-IS upstream")
	flag.StringVar(&call, "call", "", "Your callsign")
	flag.StringVar(&pass, "pass", "", "Your call pass")
	flag.StringVar(&filter, "filter", "", "Optional filter for APRS-IS server")
}

func main() {
	flag.Parse()
	if call == "" {
		fmt.Fprintf(os.Stderr, "Your callsign is required.\n")
		flag.Usage()
		os.Exit(1)
	}
	if pass == "" {
		fmt.Fprintf(os.Stderr, "Your call pass is required.\n")
		flag.Usage()
		os.Exit(1)
	}

	conn, err := textproto.Dial("tcp", server)
	if err != nil {
		log.Fatalf("Error making contact: %v", err)
	}

	if filter != "" {
		filter = fmt.Sprintf(" filter %s", filter)
	}

	conn.PrintfLine("user %s pass %s vers goaprs 0.1%s", call, pass, filter)
	for {
		line, err := conn.ReadLine()
		if err != nil {
			log.Fatalf("Error reading line:  %v", err)
		}
		if line[0] == '#' {
			log.Printf("info: %s", line)
		} else {
			msg := aprs.ParseAPRSMessage(line)
			lat, lon, err := msg.Body.Position()
			if err == nil {
				log.Printf("%s said to %s:  ``%s'' at %v,%v",
					msg.Source, msg.Dest, msg.Body, lat, lon)
			} else {
				log.Printf("%s said to %s:  ``%s''", msg.Source, msg.Dest, msg.Body)
			}
		}
	}
}
