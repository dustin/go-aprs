// Package aprsis provides an interface to APRS-IS service.
package aprsis

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/textproto"

	"github.com/dustin/go-aprs"
)

var errEmptyMsg = errors.New("empty message")
var errInvalidMsg = errors.New("invalid message")

// An APRSIS connection.
type APRSIS struct {
	conn        *textproto.Conn
	rawLog      io.Writer
	infoHandler InfoHandler
}

// InfoHandler is a handler for incoming info messages.
type InfoHandler interface {
	Info(msg string)
}

type dumbInfoHandlerT struct{}

func (d dumbInfoHandlerT) Info(msg string) {
}

var dumbInfoHandler dumbInfoHandlerT

// Next returns the next APRS message from this connection.
func (a *APRSIS) Next() (rv aprs.Frame, err error) {
	var line string
	for err == nil || err == errEmptyMsg {
		line, err = a.conn.ReadLine()
		if err != nil {
			return
		}

		fmt.Fprintf(a.rawLog, "%s\n", line)

		if len(line) > 0 && line[0] == '#' {
			a.infoHandler.Info(line)
		} else if len(line) > 0 {
			rv = aprs.ParseFrame(line)
			if !rv.IsValid() {
				err = errInvalidMsg
			}
			return rv, err
		}
	}

	return rv, errEmptyMsg
}

// SetRawLog sets a writer that will receive all raw APRS-IS messages.
func (a *APRSIS) SetRawLog(to io.Writer) {
	a.rawLog = to
}

// SetInfoHandler set a handler for APRS-IS Info messages.
func (a *APRSIS) SetInfoHandler(to InfoHandler) {
	a.infoHandler = to
}

// Dial an APRS-IS service.
func Dial(prot, addr string) (rv *APRSIS, err error) {
	var conn *textproto.Conn
	conn, err = textproto.Dial(prot, addr)
	if err != nil {
		return
	}

	return &APRSIS{conn: conn,
		rawLog:      ioutil.Discard,
		infoHandler: dumbInfoHandler,
	}, nil
}

// Close disconnects from the underlying textproto conn.
func (a *APRSIS) Close() error {
	return a.conn.Close()
}

// Send raw APRS packet using underlying textproto conn.
func (a *APRSIS) SendRawPacket(format string, args ...interface{}) error {
	return a.conn.PrintfLine(format, args...)
}

// Auth authenticates and optionally set a filter.
func (a *APRSIS) Auth(user, pass, filter string) error {
	if filter != "" {
		filter = fmt.Sprintf(" filter %s", filter)
	}

	return a.SendRawPacket("user %s pass %s vers goaprs 0.1%s",
		user, pass, filter)
}



