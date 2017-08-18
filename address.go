package aprs

import (
	"fmt"
	"strings"
)

// An Address for APRS (callsign with optional SSID)
type Address struct {
	Call string
	SSID string
}

// The string representation of an address.
func (a Address) String() string {
	rv := a.Call
	if a.SSID != "" {
		rv = fmt.Sprintf("%s-%s", a.Call, a.SSID)
	}
	return rv
}

// CallPass algorithm for APRS-IS
func (a Address) CallPass() int16 {
	rv := int16(0x73e2)
	for i := 0; i < len(a.Call); {
		rv ^= int16(a.Call[i]) << 8
		if i+1 < len(a.Call) {
			rv ^= int16(a.Call[i+1])
		}
		i += 2
	}
	return rv & 0x7fff
}

func parseAddresses(addrs []string) []Address {
	var rv []Address

	for _, s := range addrs {
		rv = append(rv, AddressFromString(s))
	}

	return rv
}

// AddressFromString builds an Addrss object from a string.
func AddressFromString(s string) Address {
	parts := strings.Split(s, "-")
	rv := Address{Call: parts[0]}
	if len(parts) > 1 {
		rv.SSID = parts[1]
	}
	return rv
}
