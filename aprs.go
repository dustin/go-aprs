package aprs

import (
	"strings"
)

type APRSMessage struct {
	Source  string
	Dest    string
	Path    []string
	Comment string
}

func ParseAPRSMessage(i string) APRSMessage {
	parts := strings.SplitN(i, ":", 2)

	srcparts := strings.SplitN(parts[0], ">", 2)
	pathparts := strings.Split(srcparts[1], ",")

	return APRSMessage{Source: srcparts[0],
		Dest: pathparts[0], Path: pathparts[1:],
		Comment: parts[1]}
}

func (m *APRSMessage) ToString() (rv string) {
	rv = strings.Join([]string{m.Source, m.Dest}, ">")
	if len(m.Path) > 0 {
		rv = strings.Join(append([]string{rv}, m.Path...), ",")
	}
	rv = strings.Join([]string{rv, m.Comment}, ":")
	return rv
}