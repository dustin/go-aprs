// Package aprs provides an Amateur Packet Radio Service messaging interface.
package aprs

import (
	"bytes"
	"strings"
)

// Info represents the information payload of an APRS packet.
type Info string

// Frame represents a complete, abstract, APRS frame.
type Frame struct {
	Original string
	Source   Address
	Dest     Address
	Path     []Address
	Body     Info
}

// IsValid is true if a message was correctly parsed.
func (d Frame) IsValid() bool {
	return d.Original != ""
}

// Type of the message.
func (b Info) Type() PacketType {
	t := PacketType(0)
	if len(b) > 0 {
		t = PacketType(b[0])
	}
	return t
}

// ParseFrame parses an APRS string into an Frame struct.
func ParseFrame(i string) Frame {
	parts := strings.SplitN(i, ":", 2)

	if len(parts) != 2 {
		return Frame{}
	}
	srcparts := strings.SplitN(parts[0], ">", 2)
	if len(srcparts) < 2 {
		return Frame{}
	}
	pathparts := strings.Split(srcparts[1], ",")

	return Frame{Original: i,
		Source: AddressFromString(srcparts[0]),
		Dest:   AddressFromString(pathparts[0]),
		Path:   parseAddresses(pathparts[1:]),
		Body:   Info(parts[1])}
}

// String forms an Frame back into its proper wire format.
func (d Frame) String() string {
	b := bytes.NewBufferString(d.Source.String())
	b.WriteByte('>')
	b.WriteString(d.Dest.String())
	for _, p := range d.Path {
		b.WriteByte(',')
		b.WriteString(p.String())
	}
	b.WriteByte(':')
	b.WriteString(string(d.Body))
	return b.String()
}
