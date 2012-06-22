package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/textproto"
	"os"

	"github.com/dustin/aprs.go"
	"github.com/dustin/aprs.go/ax25"
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

var radio io.ReadWriteCloser

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

func readNet(ch chan<- aprs.APRSMessage) {
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
		if len(line) > 0 && line[0] == '#' {
			log.Printf("info: %s", line)
		} else if len(line) > 0 {
			msg := aprs.ParseAPRSMessage(line)
			ch <- msg
		}
	}
}

func readSerial(ch chan<- aprs.APRSMessage) {
	var err error
	radio, err = rs232.OpenPort(portString, 57600, rs232.S_8N1)
	if err != nil {
		log.Fatalf("Error opening port: %s", err)
	}

	d := ax25.NewDecoder(radio)
	for {
		msg, err := d.Next()
		if err != nil {
			log.Fatalf("Error retrieving APRS message via KISS: %v", err)
		}
		ch <- msg
	}
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	if radio == nil {
		fmt.Fprintf(w, "No radio")
		return
	}

	if msg != "" {
		_, err := radio.Write([]byte{0xc0, 0x14, 0x32, 0x43, 0x00})
		if err != nil {
			log.Fatal(err)
		}
		_, err = radio.Write([]byte(msg))
		if err != nil {
			log.Fatal(err)
		}
		_, err = radio.Write([]byte{0xc0})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprintf(w, "Message sent")
	} else {
		fmt.Fprintf(w, "No message")
	}
}

func main() {
	flag.Parse()
	ch := make(chan aprs.APRSMessage)

	go reporter(ch)

	if server != "" {
		go readNet(ch)
	}

	if portString != "" {
		go readSerial(ch)
	}

	http.HandleFunc("/", sendMessage)

	log.Fatal(http.ListenAndServe(":7373", nil))
}
