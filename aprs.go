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

var setSSIDMask = byte(0x70 << 1)
var clearSSIDMask = byte(0x30 << 1)

func (a Address) kissEncode(ssidMask byte) []byte {
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

type APRSMessage struct {
	Original string
	Source   Address
	Dest     Address
	Path     []Address
	Body     MsgBody
}

func (b MsgBody) Type() PacketType {
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

func ParseAPRSMessage(i string) APRSMessage {
	parts := strings.SplitN(i, ":", 2)

	srcparts := strings.SplitN(parts[0], ">", 2)
	pathparts := strings.Split(srcparts[1], ",")

	return APRSMessage{Original: i,
		Source: AddressFromString(srcparts[0]),
		Dest:   AddressFromString(pathparts[0]),
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

func (m APRSMessage) toAX25(smask, dmask byte) []byte {
	b := &bytes.Buffer{}
	b.Write(m.Dest.kissEncode(dmask))
	b.Write(m.Source.kissEncode(smask))
	for _, p := range m.Path {
		b.Write(p.kissEncode(clearSSIDMask))
	}
	b.Write([]byte{3, 0xf0})
	b.Write([]byte(m.Body))
	return b.Bytes()
}

func (m APRSMessage) ToAX25Command() []byte {
	return m.toAX25(setSSIDMask, clearSSIDMask)
}

func (m APRSMessage) ToAX25Response() []byte {
	return m.toAX25(clearSSIDMask, setSSIDMask)
}
