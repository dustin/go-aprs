package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/textproto"
	"os"

	"github.com/dustin/aprs.go"
	"github.com/dustin/rs232.go"
)

var call, pass, filter, server, portString, rawlog string
var logWriter io.Writer = ioutil.Discard

func init() {
	flag.StringVar(&server, "server", "second.aprs.net:14580", "APRS-IS upstream")
	flag.StringVar(&portString, "port", "", "Serial port KISS thing")
	flag.StringVar(&call, "call", "", "Your callsign")
	flag.StringVar(&pass, "pass", "", "Your call pass")
	flag.StringVar(&filter, "filter", "", "Optional filter for APRS-IS server")
	flag.StringVar(&rawlog, "rawlog", "", "Path to log raw messages")
}

func reporter(ch <-chan aprs.APRSMessage) {
	for msg := range ch {
		pos, err := msg.Body.Position()
		if err == nil {
			log.Printf("%s sent a ``%v'' to %s:  ``%s'' at %v",
				msg.Source, msg.Body.Type(), msg.Dest, msg.Body, pos)
		} else {
			log.Printf("%s sent a ``%v'' to %s:  ``%s''", msg.Source,
				msg.Body.Type(), msg.Dest, msg.Body)
		}

	}
}

func processRawMessage(ch chan<- aprs.APRSMessage, line string) {
	if line[0] == '#' {
		log.Printf("info: %s", line)
	} else {
		msg := aprs.ParseAPRSMessage(line)
		ch <- msg
	}
}

func readNet(ch chan<- aprs.APRSMessage) {
	conn, err := textproto.Dial("tcp", server)
	if err != nil {
		log.Fatalf("Error making contact: %v", err)
	}

	if filter != "" {
		filter = fmt.Sprintf(" filter %s", filter)
	}

	if rawlog != "" {
		logWriter, err = os.OpenFile(rawlog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("Error opening raw log: %v", err)
		}
	}

	conn.PrintfLine("user %s pass %s vers goaprs 0.1%s", call, pass, filter)
	for {
		line, err := conn.ReadLine()
		if err != nil {
			log.Fatalf("Error reading line:  %v", err)
		}
		fmt.Fprintf(logWriter, "%s\n", line)
		processRawMessage(ch, line)
	}
}

func readSerial(ch chan<- aprs.APRSMessage) {
	port, err := rs232.OpenPort(portString, 57600, rs232.S_8N1)
	if err != nil {
		log.Fatalf("Error opening port: %s", err)
	}

	r := bufio.NewReader(&port)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Fatalf("Error reading:  %s", err)
		}
		processRawMessage(ch, string(line))
	}
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

	ch := make(chan aprs.APRSMessage)

	go reporter(ch)

	if server != "" {
		go readNet(ch)
	}

	if portString != "" {
		go readSerial(ch)
	}

	select {}

}
