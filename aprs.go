package aprs

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const coordField = `(\d{1,3})([0-5][0-9])\.(\d+)\s*([NEWS])`

var bodyRegexp = regexp.MustCompile(coordField + "/" + coordField)

var NoPositionFound = errors.New("No Positions Found")

type MsgBody string

type APRSMessage struct {
	Original string
	Source   string
	Dest     string
	Path     []string
	Body     MsgBody
}

// Get the position of the message.
func (body MsgBody) Position() (lat float64, lon float64, err error) {
	found := bodyRegexp.FindAllStringSubmatch(string(body), 8)
	// {"3722.1 N/12159.1 W", "37", "22", "1", "N", "121", "59", "1", "W"}
	if len(found) == 0 || len(found[0]) != 9 {
		return 0, 0, NoPositionFound
	}
	nums := []float64{0, 0, 0, 0}
	toparse := []string{found[0][1], found[0][2] + "." + found[0][3],
		found[0][5], found[0][6] + "." + found[0][7]}
	for i, p := range toparse {
		n, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return 0, 0, err
		}
		nums[i] = n
	}

	a := nums[0] + (nums[1] / 60)
	b := nums[2] + (nums[3] / 60)

	if found[0][4] == "S" || found[0][4] == "W" {
		a = 0 - a
	}
	if found[0][8] == "W" || found[0][8] == "S" {
		b = 0 - b
	}

	if found[0][4] == "N" || found[0][4] == "S" {
		lat = a
		lon = b
	} else {
		lat = b
		lon = a
	}

	return
}

func ParseAPRSMessage(i string) APRSMessage {
	parts := strings.SplitN(i, ":", 2)

	srcparts := strings.SplitN(parts[0], ">", 2)
	pathparts := strings.Split(srcparts[1], ",")

	return APRSMessage{Original: i,
		Source: srcparts[0],
		Dest:   pathparts[0], Path: pathparts[1:],
		Body: MsgBody(parts[1])}
}

func (m *APRSMessage) String() string {
	b := bytes.NewBufferString(m.Source)
	b.WriteByte('>')
	b.WriteString(m.Dest)
	if len(m.Path) > 0 {
		b.WriteByte(',')
		b.WriteString(strings.Join(m.Path, ","))
	}
	b.WriteByte(':')
	b.WriteString(string(m.Body))
	return b.String()
}
