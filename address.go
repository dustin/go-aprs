package aprs

import (
	"fmt"
	"strconv"
	"strings"
)

// An Address for APRS (callsign with optional SSID)
type Address struct {
	Call string
	SSID uint8
}

// The string representation of an address.
func (a Address) String() string {
	rv := a.Call
	if a.SSID != 0 {
		rv = fmt.Sprintf("%s-%d", a.Call, a.SSID)
	}
	return rv
}

// CallPass algorithm for APRS-IS
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

func parseAddresses(addrs []string) []Address {
	rv := []Address{}

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
		x, err := strconv.ParseInt(parts[1], 10, 32)
		if err == nil {
			rv.SSID = uint8(x)
		}
	}
	return rv
}
