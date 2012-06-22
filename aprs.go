package aprs

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type MsgBody string

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

type APRSMessage struct {
	Original string
	Source   Address
	Dest     Address
	Path     []Address
	Body     MsgBody
}

func (b MsgBody) Type() PacketType {
	return PacketType(b[0])
}

func parseAddress(s string) Address {
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
		rv = append(rv, parseAddress(s))
	}

	return rv
}

func ParseAPRSMessage(i string) APRSMessage {
	parts := strings.SplitN(i, ":", 2)

	srcparts := strings.SplitN(parts[0], ">", 2)
	pathparts := strings.Split(srcparts[1], ",")

	return APRSMessage{Original: i,
		Source: parseAddress(srcparts[0]),
		Dest:   parseAddress(pathparts[0]),
		Path:   parseAddresses(pathparts[1:]),
		Body:   MsgBody(parts[1])}
}

func (m *APRSMessage) String() string {
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
