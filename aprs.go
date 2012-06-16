package aprs

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const coordField = `(\d{1,3})([0-5][0-9])\.(\d+)\s*([NEWS])`
const b91chars = "[!\"#$%&'()*+,-./0123456789:;<=>?@" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`" +
	"abcdefghijklmnopqrstuvwxyz{']"

var uncompressedPositionRegexp = regexp.MustCompile(`([!=]|[/@]\d{6}[hz/])` +
	coordField + "[/DXI]" + coordField)
var compressedPositionRegexp = regexp.MustCompile("([!=/@])(" +
	b91chars + "{4})(" + b91chars + "{4})(.)(..)(.)")

var NoPositionFound = errors.New("No Positions Found")

type MsgBody string

type Position struct {
	Lat float64
	Lon float64
}

func (p Position) String() string {
	return fmt.Sprintf("{lat=%v, lon=%v}", p.Lat, p.Lon)
}

type APRSMessage struct {
	Original string
	Source   string
	Dest     string
	Path     []string
	Body     MsgBody
}

func positionUncompressed(input string) (pos Position, err error) {
	found := uncompressedPositionRegexp.FindAllStringSubmatch(input, 10)
	// {"3722.1 N/12159.1 W", "37", "22", "1", "N", "121", "59", "1", "W"}
	if len(found) == 0 || len(found[0]) != 10 {
		return pos, NoPositionFound
	}
	nums := []float64{0, 0, 0, 0}
	toparse := []string{found[0][2], found[0][3] + "." + found[0][4],
		found[0][6], found[0][7] + "." + found[0][8]}
	for i, p := range toparse {
		n, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return pos, err
		}
		nums[i] = n
	}

	a := nums[0] + (nums[1] / 60)
	b := nums[2] + (nums[3] / 60)

	if found[0][5] == "S" || found[0][5] == "W" {
		a = 0 - a
	}
	if found[0][9] == "W" || found[0][9] == "S" {
		b = 0 - b
	}

	if found[0][5] == "N" || found[0][5] == "S" {
		pos.Lat = a
		pos.Lon = b
	} else {
		pos.Lat = b
		pos.Lon = a
	}

	// log.Printf("uncomp matched %#v -> %v,%v", found, lat, lon)

	return
}

func decodeBase91(s []byte) int {
	if len(s) != 4 {
		return 0
	}
	return ((int(s[0]) - 33) * 91 * 91 * 91) + ((int(s[1]) - 33) * 91 * 91) +
		(int(s[2]-33) * 91) + int(s[3]) - 33
}

func positionCompressed(input string) (pos Position, err error) {
	found := compressedPositionRegexp.FindAllStringSubmatch(input, 10)
	// {"/]\"4-}Foo !w6", "/", "]\"4-", "}Foo", " ", "!w", "6"}}
	if len(found) == 0 || len(found[0]) != 7 {
		return pos, NoPositionFound
	}

	// Lat = 90 - ((y1-33) x 91^3 + (y2-33) x 91^2 + (y3-33) x 91 + y4-33) / 380926
	// Long = -180 + ((x1-33) x 91^3 + (x2-33) x 91^2 + (x3-33) x 91 + x4-33) / 190463

	pos.Lat = 90 - float64(decodeBase91([]byte(found[0][2])))/380926
	pos.Lon = -180 + float64(decodeBase91([]byte(found[0][3])))/190463

	// log.Printf("comp matched %#v (%v)-> %v,%v", found, found[0][4], lat, lon)

	return pos, nil
}

// Get the position of the message.
func (body MsgBody) Position() (pos Position, err error) {
	pos, err = positionUncompressed(string(body))
	if err == nil {
		return
	}
	pos, err = positionCompressed(string(body))
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
