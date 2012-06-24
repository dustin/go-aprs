package aprs

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Info string

type Address struct {
	Call string
	SSID int
}

func (a Address) String() string {
	rv := a.Call
	if a.SSID != 0 {
		rv = fmt.Sprintf("%s-%d", a.Call, a.SSID)
	}
	return rv
}

func (a Address) CallPass() (rv int16) {
	rv = 0x73e2
	for i := 0; i < len(a.Call); {
		rv ^= int16(a.Call[i]) << 8
		rv ^= int16(a.Call[i+1])
		i += 2
	}
	rv &= 0x7fff
	return
}

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

func AddressFromString(s string) Address {
	parts := strings.Split(s, "-")
	rv := Address{Call: parts[0]}
	if len(parts) > 1 {
		x, err := strconv.ParseInt(parts[1], 10, 32)
		if err == nil {
			rv.SSID = int(x)
		}
	}
	return rv
}

func parseAddresses(addrs []string) []Address {
	rv := []Address{}

	for _, s := range addrs {
		rv = append(rv, AddressFromString(s))
	}

	return rv
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

func (m *APRSData) String() string {
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
