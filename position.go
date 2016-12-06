package aprs

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
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

// ErrNoPosition is returned when no geo positions could be found in a message.
var ErrNoPosition = errors.New("no positions found")

// ErrTruncatedMsg is returned when a message is incomplete.
var ErrTruncatedMsg = errors.New("truncated message")

// Symbol represents the map marker symbol for an object or station.
type Symbol struct {
	Table  byte
	Symbol byte
}

// IsPrimary is true if this symbol is part of the primary symbol table.
func (s Symbol) IsPrimary() bool {
	return s.Table != '\\'
}

// Name is the name of the symbol.
func (s Symbol) Name() (rv string) {
	m := primarySymbolMap
	if !s.IsPrimary() {
		m = alternateSymbolMap
	}
	return m[s.Symbol]
}

// Glyph returns a textual representation of this Symbol.
func (s Symbol) Glyph() string {
	return symbolGlyphs[s.Name()]
}

func (s Symbol) String() (rv string) {
	g := s.Glyph()
	if g == "" {
		rv = fmt.Sprintf("{%c%c: %s}", s.Table, s.Symbol, s.Name())
	} else {
		rv = fmt.Sprintf("{%c%c: %s - %s}", s.Table, s.Symbol, s.Name(), g)
	}
	return
}

// Velocity represents the course and speed of an object or station.
type Velocity struct {
	Course float64
	Speed  float64
}

// Position contains all of the information necessary for placing an object on a map.
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

func uncompressedParser(input string) (pos Position, err error) {
	// lat:8 symtab:1 lon:9 sym:1
	if len(input) < 19 {
		return pos, ErrTruncatedMsg
	}

	pos.Symbol.Table = input[8]
	pos.Symbol.Symbol = input[18]

	nums := []float64{0, 0, 0, 0}
	toparse := []string{input[0:2], input[2:7], input[9:12], input[12:17]}

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
		return pos, fmt.Errorf("invalid position ambiguity %d from %v",
			pos.Ambiguity, input)
	}
	if offby > 0 {
		a += offby
		b += offby
	}

	if input[7] == 'S' {
		a = 0 - a
	}
	if input[17] == 'W' {
		b = 0 - b
	}

	pos.Lat = a
	pos.Lon = b

	ext := input[19:]
	if len(ext) >= 7 && pos.Symbol.Symbol != '_' && ext[3] == '/' {
		fmt.Sscanf(ext, "%f/%f",
			&pos.Velocity.Course, &pos.Velocity.Speed)
		pos.Velocity.Speed *= 1.852
	}

	return
}

func positionUncompressed(input string) (Position, error) {
	found := uncompressedPositionRegexp.FindAllStringSubmatch(input, 10)
	// {"=3722.1 N/12159.1 W-", "=", "37", "22", "1 ", "N", "/", "121", "59", "1 ", "W", "-", ""}
	if len(found) == 0 || len(found[0]) != 13 {
		return Position{}, ErrNoPosition
	}
	pos := Position{}
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
		return pos, fmt.Errorf("invalid position ambiguity %d from %v",
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

	return pos, nil
}

func decodeBase91(s []byte) int {
	if len(s) != 4 {
		return 0
	}
	return ((int(s[0]) - 33) * 91 * 91 * 91) + ((int(s[1]) - 33) * 91 * 91) +
		(int(s[2]-33) * 91) + int(s[3]) - 33
}

func positionCompressed(input string) (Position, error) {
	found := compressedPositionRegexp.FindAllStringSubmatch(input, 10)
	// {"/]\"4-}Foo !w6", "/", "]\"4-", "}Foo", " ", "!w", "6"}}
	if len(found) == 0 || len(found[0]) != 7 {
		return Position{}, ErrNoPosition
	}

	// Lat = 90 - ((y1-33) x 91^3 + (y2-33) x 91^2 + (y3-33) x 91 + y4-33) / 380926
	// Long = -180 + ((x1-33) x 91^3 + (x2-33) x 91^2 + (x3-33) x 91 + x4-33) / 190463

	pos := Position{
		Lat: 90 - float64(decodeBase91([]byte(found[0][2])))/380926,
		Lon: -180 + float64(decodeBase91([]byte(found[0][3])))/190463,
	}

	cs := found[0][5]
	if cs[0] != ' ' && cs[1] != ' ' && int(cs[0]) >= '!' && int(cs[0]) <= 'z' {
		pos.Velocity.Course = (float64(cs[0]) - 33) * 4
		if pos.Velocity.Course == 0 {
			pos.Velocity.Course = 360
		}
		pos.Velocity.Speed = 1.852 * (math.Pow(1.08, float64(cs[1]-33)) - 1)

	}

	return pos, nil
}

func positionOld(t string) (Position, error) {
	pos, err := positionUncompressed(t)
	if err == nil {
		return pos, err
	}
	return positionCompressed(t)
}

func compressedParser(input string) (Position, error) {
	if len(input) < 12 {
		return Position{}, ErrTruncatedMsg
	}
	pos := Position{}
	pos.Symbol.Table = input[0]
	pos.Symbol.Symbol = input[9]
	pos.Lat = 90 - float64(decodeBase91([]byte(input[1:5])))/380926
	pos.Lon = -180 + float64(decodeBase91([]byte(input[5:9])))/190463
	if input[10] != ' ' && input[11] != ' ' && int(input[10]) >= '!' && int(input[10]) <= 'z' {
		pos.Velocity.Course = (float64(input[10]) - 33) * 4
		if pos.Velocity.Course == 0 {
			pos.Velocity.Course = 360
		}
		pos.Velocity.Speed = 1.852 * (math.Pow(1.08, float64(input[11]-33)) - 1)
	}
	return pos, nil
}

func newParser(input string, uncompressed bool) (Position, error) {
	if uncompressed {
		return uncompressedParser(input)
	}
	return compressedParser(input)
}

// Position gets the position of the message.
func (body Info) Position() (Position, error) {
	switch body.Type() {
	case '!', '=':
		t := string(body)
		return newParser(t[1:], unicode.IsDigit(rune(t[1])))
	case '/', '@':
		t := string(body[8:])
		return newParser(t, unicode.IsDigit(rune(body[8])))
	case ';':
		t := string(body)
		if len(t) < 19 {
			return Position{}, ErrTruncatedMsg
		}
		return newParser(t[18:], unicode.IsDigit(rune(body[18])))
		// t := string(body[1:])
		// name := strings.TrimSpace(t[1:9])
		// live := t[9] == '*'
		// ts := t[10:17]
	case ')':
		// item
	}
	return positionOld(string(body))
}
