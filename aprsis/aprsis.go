package aprsis

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/textproto"

	"github.com/dustin/go-aprs"
)

var emptyMessage = errors.New("empty message")

type APRSIS struct {
	conn        *textproto.Conn
	rawLog      io.Writer
	infoHandler InfoHandler
}

type InfoHandler interface {
	Info(msg string)
}

type dumbInfoHandlerT struct{}

func (d dumbInfoHandlerT) Info(msg string) {
}

var dumbInfoHandler dumbInfoHandlerT

func (a *APRSIS) Next() (rv aprs.APRSMessage, err error) {
	var line string
	for err == nil || err == emptyMessage {
		line, err = a.conn.ReadLine()
		if err != nil {
			return
		}

		fmt.Fprintf(a.rawLog, "%s\n", line)

		if len(line) > 0 && line[0] == '#' {
			a.infoHandler.Info(line)
		} else if len(line) > 0 {
			rv = aprs.ParseAPRSMessage(line)
			return rv, nil
		}
	}

	return rv, emptyMessage
}

func (a *APRSIS) SetRawLog(to io.Writer) {
	a.rawLog = to
}

func (a *APRSIS) SetInfoHandler(to InfoHandler) {
	a.infoHandler = to
}

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

func (a *APRSIS) Auth(user, pass, filter string) error {
	if filter != "" {
		filter = fmt.Sprintf(" filter %s", filter)
	}

	return a.conn.PrintfLine("user %s pass %s vers goaprs 0.1%s",
		user, pass, filter)
}
