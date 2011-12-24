package aprs

import (
	"strings"
)

type APRSMessage struct {
	Original string
	Source   string
	Dest     string
	Path     []string
	Body     string
}

func ParseAPRSMessage(i string) APRSMessage {
	parts := strings.SplitN(i, ":", 2)

	srcparts := strings.SplitN(parts[0], ">", 2)
	pathparts := strings.Split(srcparts[1], ",")

	return APRSMessage{Original: i,
		Source: srcparts[0],
		Dest: pathparts[0], Path: pathparts[1:],
		Body: parts[1]}
}

func (m *APRSMessage) ToString() (rv string) {
	rv = strings.Join([]string{m.Source, m.Dest}, ">")
	if len(m.Path) > 0 {
		rv = strings.Join(append([]string{rv}, m.Path...), ",")
	}
	rv = strings.Join([]string{rv, m.Body}, ":")
	return rv
}