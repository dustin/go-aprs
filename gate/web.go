package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dustin/go-aprs"
	"github.com/dustin/go-aprs/ax25"
)

func init() {
	http.HandleFunc("/", sendMessage)
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

		msg := aprs.APRSData{
			Source: aprs.AddressFromString(src),
			Dest:   aprs.AddressFromString(dest),
			Path: []aprs.Address{
				aprs.AddressFromString("WIDE2-2")},
			Body: aprs.Info(text),
		}

		body := ax25.EncodeAPRSCommand(msg)
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
