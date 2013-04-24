// AX.25 encoding and decoding lib.
package ax25

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/dustin/go-aprs"
)

const reasonableSize = 14

var shortMessage = errors.New("short message")
var truncatedMessage = errors.New("truncated message")

var setSSIDMask = byte(0x70 << 1)
var clearSSIDMask = byte(0x30 << 1)

func parseAddr(in []byte) aprs.Address {
	out := make([]byte, len(in))
	for i, b := range in {
		out[i] = b >> 1
	}
	rv := aprs.Address{
		Call: strings.TrimSpace(string(out[:len(out)-1])),
		SSID: uint8((out[len(out)-1]) & 0xf),
	}
	return rv
}

func decodeMessage(frame []byte) (rv aprs.APRSData, err error) {
	frame = frame[:len(frame)-1]

	if len(frame) < reasonableSize {
		err = shortMessage
		return
	}

	rv.Source = parseAddr(frame[8:15])
	rv.Dest = parseAddr(frame[1:8])

	rv.Path = []aprs.Address{}

	frame = frame[15:]
	for len(frame) > 7 && frame[0] != 3 {
		rv.Path = append(rv.Path, parseAddr(frame[:7]))
		frame = frame[7:]
	}

	if len(frame) < 2 || frame[0] != 3 || frame[1] != 0xf0 {
		err = truncatedMessage
		return
	}

	rv.Body = aprs.Info(string(frame[2:]))

	return
}

// An AX.25 message decoder.
type Decoder struct {
	r *bufio.Reader
}

// Get the next message.
func (d *Decoder) Next() (aprs.APRSData, error) {
	frame := []byte{}
	var err error
	for len(frame) < reasonableSize {
		frame, err = d.r.ReadBytes(byte(0xc0))
		if err != nil {
			return aprs.APRSData{}, err
		}
	}
	return decodeMessage(frame)
}

// Get a new decoder over this reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{bufio.NewReader(r)}
}

func addressEncode(a aprs.Address, ssidMask byte) []byte {
	rv := make([]byte, 7)
	for i := 0; i < len(rv); i++ {
		rv[i] = ' '
	}
	for i, c := range a.Call {
		rv[i] = byte(c) << 1
	}
	rv[6] = ssidMask | (byte(a.SSID) << 1)
	return rv
}

func toAX25(m aprs.APRSData, smask, dmask byte) []byte {
	b := &bytes.Buffer{}
	b.Write(addressEncode(m.Dest, dmask))
	mask := smask
	if len(m.Path) == 0 {
		mask |= 1
	}
	b.Write(addressEncode(m.Source, smask))
	for i, p := range m.Path {
		mask = clearSSIDMask
		if i == len(m.Path)-1 {
			mask |= 1
		}
		b.Write(addressEncode(p, mask))
	}
	b.Write([]byte{3, 0xf0})
	b.Write([]byte(m.Body))
	return b.Bytes()
}

// Encode an APRS command to an AX.25 frame.
func EncodeAPRSCommand(m aprs.APRSData) []byte {
	return toAX25(m, setSSIDMask, clearSSIDMask)
}

// Encode an APRS response to an AX.25 frame.
func EncodeAPRSResponse(m aprs.APRSData) []byte {
	return toAX25(m, clearSSIDMask, setSSIDMask)
}
