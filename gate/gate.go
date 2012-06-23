package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dustin/go-aprs"
	"github.com/dustin/go-aprs/aprsis"
	"github.com/dustin/go-aprs/ax25"
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

type loggingInfoHandler struct{}

func (*loggingInfoHandler) Info(msg string) {
	log.Printf("info: %s", msg)

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

	is, err := aprsis.Dial("tcp", server)
	if err != nil {
		log.Fatalf("Error making contact: %v", err)
	}

	if rawlog != "" {
		logWriter, err := os.OpenFile(rawlog,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("Error opening raw log: %v", err)
		}
		is.SetRawLog(logWriter)
	}

	is.SetInfoHandler(&loggingInfoHandler{})

	is.Auth(call, pass, filter)

	for {
		msg, err := is.Next()
		if err != nil {
			log.Fatalf("Error reading line:  %v", err)
		}
		ch <- msg
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

func sendMessage(rw http.ResponseWriter, r *http.Request) {
	src := r.FormValue("src")
	dest := r.FormValue("dest")
	text := r.FormValue("msg")
	if radio == nil {
		fmt.Fprintf(rw, "No radio")
		return
	}

	if text != "" {
		d := hex.Dumper(os.Stdout)
		defer d.Close()
		w := io.MultiWriter(d, radio)

		n, err := w.Write([]byte{0xc0, 0x00})
		if err != nil {
			log.Fatal(err)
		}
		if n != 2 {
			log.Fatalf("Expected to write two bytes, wrote %v", n)
		}

		msg := aprs.APRSMessage{
			Source: aprs.AddressFromString(src),
			Dest:   aprs.AddressFromString(dest),
			Path: []aprs.Address{
				aprs.AddressFromString("WIDE2-2")},
			Body: aprs.MsgBody(text),
		}

		body := msg.ToAX25Command()
		n, err = w.Write(body)
		if err != nil {
			log.Fatal(err)
		}
		if n != len(body) {
			log.Fatalf("Expected to write %v bytes, wrote %v", len(body), n)
		}

		n, err = w.Write([]byte{0xc0})
		if err != nil {
			log.Fatal(err)
		}
		if n != 1 {
			log.Fatalf("Expected to write 1 byte, wrote %v", n)
		}

		fmt.Fprintf(rw, "Message sent")
	} else {
		fmt.Fprintf(rw, "No message")
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
