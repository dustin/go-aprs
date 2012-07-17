package main

import (
	"fmt"
	"log"
	"math"
	"net"

	"github.com/dustin/go-aprs"
)

// Server filters limit the packets that are received over APRS-IS
type Filter interface {
	Matches(d aprs.APRSData) bool
}

// A filter made up of other filters
type CompositeFilter struct {
	Positive []Filter
	Negative []Filter
}

func (c *CompositeFilter) Matches(d aprs.APRSData) bool {
	rv := false
	for _, f := range c.Positive {
		if f.Matches(d) {
			rv = true
			break
		}
	}
	if rv {
		for _, f := range c.Negative {
			if f.Matches(d) {
				rv = false
				break
			}
		}
	}
	return rv
}

type Point struct {
	Lat float64
	Lon float64
}

func d2r(d float64) float64 {
	return d * 0.0174532925
}

func (p Point) RadLat() float64 {
	return d2r(p.Lat)
}

func (p Point) RadLon() float64 {
	return d2r(p.Lon)
}

func (p Point) Distance(p2 Point) float64 {
	r := 6371.01
	return math.Acos((math.Sin(p.RadLat())*
		math.Sin(p2.RadLat()))+
		(math.Cos(p.RadLat())*math.Cos(p2.RadLat())*
			math.Cos(p.RadLon()-p2.RadLon()))) * r
}

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
	ln, err := net.Listen(n, addr)
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
