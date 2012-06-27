// Amateur Packet Radio Service library.
package aprs

import (
	"bytes"
	"strings"
)

type Info string

type APRSData struct {
	Original string
	Source   Address
	Dest     Address
	Path     []Address
	Body     Info
}

func (d APRSData) IsValid() bool {
	return d.Original != ""
}

func (b Info) Type() PacketType {
	t := byte(0)
	if len(b) > 0 {
		t = b[0]
	}
	return PacketType(t)
}

func ParseAPRSData(i string) APRSData {
	parts := strings.SplitN(i, ":", 2)

	if len(parts) != 2 {
		return APRSData{}
	}
	srcparts := strings.SplitN(parts[0], ">", 2)
	if len(srcparts) < 2 {
		return APRSData{}
	}
	pathparts := strings.Split(srcparts[1], ",")

	return APRSData{Original: i,
		Source: AddressFromString(srcparts[0]),
		Dest:   AddressFromString(pathparts[0]),
		Path:   parseAddresses(pathparts[1:]),
		Body:   Info(parts[1])}
}

func (m APRSData) String() string {
	b := bytes.NewBufferString(m.Source.String())
	b.WriteByte('>')
	b.WriteString(m.Dest.String())
	for _, p := range m.Path {
		b.WriteByte(',')
		b.WriteString(p.String())
	}
	b.WriteByte(':')
	b.WriteString(string(m.Body))
	return b.String()
}
