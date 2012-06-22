package ax25

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/dustin/go-aprs"
)

const reasonableSize = 14

var shortMessage = errors.New("short message")
var truncatedMessage = errors.New("truncated message")

func parseAddr(in []byte) aprs.Address {
	out := make([]byte, len(in))
	for i, b := range in {
		out[i] = b >> 1
	}
	rv := aprs.Address{
		Call: strings.TrimSpace(string(out[:len(out)-1])),
		SSID: int((out[len(out)-1]) & 0xf),
	}
	return rv
}

func decodeMessage(frame []byte) (rv aprs.APRSMessage, err error) {
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

	rv.Body = aprs.MsgBody(string(frame[2:]))

	return
}

// An AX.25 message decoder.
type Decoder struct {
	r *bufio.Reader
}

// Get the next message.
func (d *Decoder) Next() (aprs.APRSMessage, error) {
	frame := []byte{}
	var err error
	for len(frame) < reasonableSize {
		frame, err = d.r.ReadBytes(byte(0xc0))
		if err != nil {
			return aprs.APRSMessage{}, err
		}
	}
	return decodeMessage(frame)
}

// Get a new decoder over this reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{bufio.NewReader(r)}
}
