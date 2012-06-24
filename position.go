package aprs

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const coordField = `(\d{1,3})([0-5 ][0-9 ])\.([0-9 ]+)([NEWS])`
const b91chars = "[!\"#$%&'()*+,-./0123456789:;<=>?@" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`" +
	"abcdefghijklmnopqrstuvwxyz{']"

const symbolTables = `[0-9/\\A-z]`

var uncompressedPositionRegexp = regexp.MustCompile(`([!=]|[/@\*]\d{6}[hz/])` +
	coordField + "(" + symbolTables + ")" + coordField + "(.)([0-3][0-9]{2}/[0-9]{3})?")
var compressedPositionRegexp = regexp.MustCompile("([!=/@])(" +
	b91chars + "{4})(" + b91chars + "{4})(.)(..)(.)")

var NoPositionFound = errors.New("No Positions Found")

type Symbol struct {
	Table  byte
	Symbol byte
}

func (s Symbol) IsPrimary() bool {
	return s.Table != '\\'
}

func (s Symbol) Name() (rv string) {
	m := primarySymbolMap
	if !s.IsPrimary() {
		m = alternateSymbolMap
	}
	return m[s.Symbol]
}

func (s Symbol) Glyph() string {
	return symbolGlyphs[s.Name()]
}

func (s Symbol) String() (rv string) {
	g := s.Glyph()
	if g == "" {
		rv = fmt.Sprintf("{%c%c: %s}", s.Table, s.Symbol, s.Name())
	} else {
		rv = fmt.Sprintf("{%c%c: %s -  %s}", s.Table, s.Symbol, s.Name(), g)
	}
	return
}

type Velocity struct {
	Course float64
	Speed  float64
}

type Position struct {
	Lat       float64
	Lon       float64
	Ambiguity int
	Velocity  Velocity
	Symbol    Symbol
}

func (p Position) String() string {
	return fmt.Sprintf("{lat=%v, lon=%v, amb=%v, sym=%v}",
		p.Lat, p.Lon, p.Ambiguity, p.Symbol)
}

func positionUncompressed(input string) (pos Position, err error) {
	found := uncompressedPositionRegexp.FindAllStringSubmatch(input, 10)
	// {"=3722.1 N/12159.1 W-", "=", "37", "22", "1 ", "N", "/", "121", "59", "1 ", "W", "-", ""}
	if len(found) == 0 || len(found[0]) != 13 {
		return pos, NoPositionFound
	}
	pos.Symbol.Table = found[0][6][0]
	pos.Symbol.Symbol = found[0][11][0]
	nums := []float64{0, 0, 0, 0}
	toparse := []string{found[0][2], found[0][3] + "." + found[0][4],
		found[0][7], found[0][8] + "." + found[0][9]}
	for i, p := range toparse {
		converted := strings.Map(func(r rune) (rv rune) {
			rv = r
			if r == ' ' {
				pos.Ambiguity++
				rv = '0'
			}
			return
		}, p)
		n, err := strconv.ParseFloat(converted, 64)
		if err != nil {
			return pos, err
		}
		nums[i] = n
	}

	a := nums[0] + (nums[1] / 60)
	b := nums[2] + (nums[3] / 60)

	pos.Ambiguity /= 2
	offby := 0.0
	switch pos.Ambiguity {
	case 0:
		// This is exact
	case 1:
		// Nearest 1/10 of a minute
		offby = 0.05 / 60.0
	case 2:
		// Nearest minute
		offby = 0.5 / 60.0
	case 3:
		// Nearest 10 minutes
		offby = 5.0 / 60.0
	case 4:
		// Nearest degree
		offby = 0.5
	default:
		return pos, fmt.Errorf("Invalid position ambiguity %d from %v",
			pos.Ambiguity, found[0][0])
	}
	if offby > 0 {
		a += offby
		b += offby
	}

	if found[0][5] == "S" || found[0][5] == "W" {
		a = 0 - a
	}
	if found[0][10] == "W" || found[0][10] == "S" {
		b = 0 - b
	}

	if found[0][5] == "N" || found[0][5] == "S" {
		pos.Lat = a
		pos.Lon = b
	} else {
		pos.Lat = b
		pos.Lon = a
	}

	if found[0][12] != "" && pos.Symbol.Symbol != '_' {
		fmt.Sscanf(found[0][12], "%f/%f",
			&pos.Velocity.Course, &pos.Velocity.Speed)
		pos.Velocity.Speed *= 1.852
	}

	// log.Printf("uncomp matched %#v -> %v", found, pos)

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
func (body Info) Position() (pos Position, err error) {
	pos, err = positionUncompressed(string(body))
	if err == nil {
		return
	}
	pos, err = positionCompressed(string(body))
	return
}
