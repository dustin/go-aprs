package aprs

import (
	"bytes"
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

func (m *APRSMessage) ToString() string {
	b := bytes.NewBufferString(m.Source)
	b.WriteByte('>')
	b.WriteString(m.Dest)
	if len(m.Path) > 0 {
		b.WriteByte(',')
		b.WriteString(strings.Join(m.Path, ","))
	}
	b.WriteByte(':')
	b.WriteString(m.Body)
	return b.String()
}