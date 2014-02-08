package ax25

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/dustin/go-aprs"
)

const (
	christmasMsg = "KG6HWF>APX200,WIDE1-1,WIDE2-1:=3722.1 N/12159.1 W-Merry Christmas!"
)

func TestKISS(t *testing.T) {
	v := aprs.ParseAPRSData(christmasMsg)
	bc := EncodeAPRSCommand(v)
	t.Logf("Command:\n" + hex.Dump(bc))

	br := EncodeAPRSResponse(v)
	t.Logf("Response:\n" + hex.Dump(br))
}

func TestAddressConversion(t *testing.T) {
	testaddrs := []struct {
		Src     string
		AX25Cmd []byte
		AX25Res []byte
	}{
		{"KG6HWF",
			[]byte{0x96, 0x8e, 0x6c, 0x90, 0xae, 0x8c, 0xe0},
			[]byte{0x96, 0x8e, 0x6c, 0x90, 0xae, 0x8c, 0x60}},
		{"KG6HWF-9",
			[]byte{0x96, 0x8e, 0x6c, 0x90, 0xae, 0x8c, 0xf2},
			[]byte{0x96, 0x8e, 0x6c, 0x90, 0xae, 0x8c, 0x72}},
	}

	for _, ta := range testaddrs {
		a := aprs.AddressFromString(ta.Src)
		a25c := addressEncode(a, setSSIDMask)
		if !reflect.DeepEqual(a25c, ta.AX25Cmd) {
			t.Fatalf("Expected %v for AX25d %v, got %v",
				ta.AX25Cmd, ta.Src, a25c)
		}
		a25r := addressEncode(a, clearSSIDMask)
		if !reflect.DeepEqual(a25r, ta.AX25Res) {
			t.Fatalf("Expected %v for AX25d %v, got %v",
				ta.AX25Res, ta.Src, a25r)
		}
	}
}
