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

func sendMessage(w http.ResponseWriter, r *http.Request) {
	src := r.FormValue("src")
	dest := r.FormValue("dest")
	text := r.FormValue("msg")
	if radio == nil {
		http.Error(w, "No radio", 500)
		return
	}

	if text != "" {
		d := hex.Dumper(os.Stdout)
		defer d.Close()
		mw := io.MultiWriter(d, radio)

		_, err := mw.Write([]byte{0xc0, 0x00})
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Printf("Error writing command: %v", err)
			return
		}

		msg := aprs.Frame{
			Source: aprs.AddressFromString(src),
			Dest:   aprs.AddressFromString(dest),
			Path: []aprs.Address{
				aprs.AddressFromString("WIDE2-2")},
			Body: aprs.Info(text),
		}

		body := ax25.EncodeAPRSCommand(msg)
		_, err = mw.Write(body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Printf("Error writing command: %v", err)
			return
		}

		_, err = mw.Write([]byte{0xc0})
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Printf("Error finishing command: %v", err)
			return
		}

		fmt.Fprintf(w, "Message sent")
	} else {
		http.Error(w, "No message", 400)
	}
}
