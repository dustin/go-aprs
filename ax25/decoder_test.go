package ax25

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/dustin/aprs.go"
)

func TestCapture(t *testing.T) {
	f, err := os.Open("radio.sample")
	if err != nil {
		t.Fatalf("Error opening sample file: %v", err)
	}
	defer f.Close()

	expected := []aprs.APRSMessage{
		aprs.APRSMessage{Source: aprs.Address{Call: "N6WKZ", SSID: 3},
			Dest: aprs.Address{Call: "APU25N", SSID: 0},
			Path: []aprs.Address{aprs.Address{Call: "WR6ABD", SSID: 0}},
			Body: "=3746.42N112226.00W# {UIV32N}\r"},
		aprs.APRSMessage{Source: aprs.Address{Call: "W1EJ", SSID: 10},
			Dest: aprs.Address{Call: "APT311", SSID: 0},
			Path: []aprs.Address{aprs.Address{Call: "WB6TMS", SSID: 5},
				aprs.Address{Call: "N6ZX", SSID: 3},
				aprs.Address{Call: "WIDE2", SSID: 0}},
			Body: "/210725z3814.29N/12236.93W>275/000/A=000013/ED J SAG"},
		aprs.APRSMessage{Source: aprs.Address{Call: "WR6ABD", SSID: 0},
			Dest: aprs.Address{Call: "APN382", SSID: 0},
			Path: []aprs.Address{},
			Body: "!3706.66NS12150.69W#PHG5730 W1,NCAn Loma Prieta LPRC.net A=003980\r"},
		aprs.APRSMessage{Source: aprs.Address{Call: "N6ACK", SSID: 1},
			Dest: aprs.Address{Call: "APRS", SSID: 0},
			Path: []aprs.Address{},
			Body: "}WR6ABD>APN382,TCPIP*,N6ACK-1*:!3706.66NS12150.69W#PHG5730 W1,NCAn Loma Prieta LPRC.net A=003980"},
		aprs.APRSMessage{Source: aprs.Address{Call: "N6ACK", SSID: 1},
			Dest: aprs.Address{Call: "APRS", SSID: 0},
			Path: []aprs.Address{},
			Body: "}AC6SL-4>APD225,TCPIP*,N6ACK-1*:!3707.94NI12207.23W& receive-only-aprsd"},
		aprs.APRSMessage{Source: aprs.Address{Call: "CARSON", SSID: 0},
			Dest: aprs.Address{Call: "APN391", SSID: 0},
			Path: []aprs.Address{aprs.Address{Call: "ECHO", SSID: 0},
				aprs.Address{Call: "N6ZX", SSID: 3},
				aprs.Address{Call: "WIDE2", SSID: 0}},
			Body: "!3841.68N111959.36W#PHG7636/NCAn,TEMPn/WG6D/Carson Pass, CA/A=008573\r"},
		aprs.APRSMessage{Source: aprs.Address{Call: "KE6KYI", SSID: 0},
			Dest: aprs.Address{Call: "APU25N", SSID: 0},
			Path: []aprs.Address{aprs.Address{Call: "K6TUO", SSID: 3},
				aprs.Address{Call: "N6ZX", SSID: 3},
				aprs.Address{Call: "WIDE2", SSID: 0}},
			Body: "@210726z3751.53N/12012.83W_213/000g000t063r000p000P000h45b10096APRS/CWOP Weather\r"},
		aprs.APRSMessage{Source: aprs.Address{Call: "N6ACK", SSID: 1},
			Dest: aprs.Address{Call: "APRS", SSID: 0},
			Path: []aprs.Address{},
			Body: "}N6VIG-9>SW4QTY,TCPIP*,N6ACK-1*:`1Q\x1el ?>\\\"4m}"},
		aprs.APRSMessage{Source: aprs.Address{Call: "KG6ZLQ", SSID: 12},
			Dest: aprs.Address{Call: "3X5SRR", SSID: 0},
			Path: []aprs.Address{aprs.Address{Call: "ECHO", SSID: 0},
				aprs.Address{Call: "WIDE1", SSID: 0},
				aprs.Address{Call: "N6ZX", SSID: 3},
				aprs.Address{Call: "WIDE2", SSID: 0}},
			Body: "`0Z)l\"{j/\"IN}"},
		aprs.APRSMessage{Source: aprs.Address{Call: "N6ACK", SSID: 1},
			Dest: aprs.Address{Call: "APRS", SSID: 0},
			Path: []aprs.Address{},
			Body: "}WA6BAY-1>APRS,TCPIP*,N6ACK-1*:!!0000008101F905B0276B02E803E8----00AC001A00000000"},
		aprs.APRSMessage{Source: aprs.Address{Call: "W6SIG", SSID: 0},
			Dest: aprs.Address{Call: "APS228", SSID: 0},
			Path: []aprs.Address{aprs.Address{Call: "W6CX", SSID: 3}},
			Body: "=3834.22N/12118.36WoPHG33D0 CalEMA-Mather\r"},
		aprs.APRSMessage{Source: aprs.Address{Call: "KI6ASH", SSID: 0},
			Dest: aprs.Address{Call: "S7SXWV", SSID: 0},
			Path: []aprs.Address{aprs.Address{Call: "WA6TOW", SSID: 2},
				aprs.Address{Call: "W6CX", SSID: 3},
				aprs.Address{Call: "WIDE2", SSID: 0}},
			Body: "`24gl \x1c>/'\"3u}MT-RTG|%V%`'n|!wwU!|3"}}

	got := []aprs.APRSMessage{}

	d := NewDecoder(f)
	for err == nil {
		var m aprs.APRSMessage
		m, err = d.Next()
		if err == nil {
			got = append(got, m)
		}
	}
	if err != io.EOF {
		t.Fatalf("Error reading stream: %v", err)
	}

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("Expected:\n%#v\nGot:\n%#v", expected, got)
	}
}
