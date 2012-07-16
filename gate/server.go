package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dustin/go-aprs"
)

func handleIS(conn net.Conn, b *broadcaster) {
	ch := make(chan aprs.APRSData, 100)

	_, err := fmt.Fprintf(conn, "# goaprs\n")
	if err != nil {
		log.Printf("Error sending banner: %v", err)
	}

	b.Register(ch)
	defer b.Unregister(ch)

	for m := range ch {
		_, err = conn.Write([]byte(m.String() + "\n"))
		if err != nil {
			log.Printf("Error on connection:  %v", err)
			return
		}
	}
}

func startIS(n, addr string, b *broadcaster) {
	ln, err := net.Listen("tcp", ":10152")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connections: %v", err)
			continue
		}
		go handleIS(conn, b)
	}
}
