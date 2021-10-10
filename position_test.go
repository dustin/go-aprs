package aprs

import (
	"fmt"
	"testing"
)

func TestSymbol(t *testing.T) {
	tests := []struct {
		in Symbol
		s  string
	}{
		{Symbol{'/', '\''}, `{/': Plane sm - ✈}`},
		{Symbol{'\\', '\''}, `{\': Crash site}`},
		{Symbol{'/', 'a'}, "{/a: Ambulance - \u2620}"},
	}

	for _, test := range tests {
		if test.in.String() != test.s {
			t.Errorf("On %#v.String() = %v, want %v", test.in, test.in.String(), test.s)
		}
	}
}

func TestPosition(t *testing.T) {
	p := Position{37, -121, 2, Velocity{15, 31}, Symbol{'/', 'a'}}
	exp := "{lat=37, lon=-121, amb=2, sym={/a: Ambulance - \u2620}}"
	if p.String() != exp {
		t.Errorf("for %#v, got %v, want %v", p, p, exp)
	}
}

func TestDecodeBase91(t *testing.T) {
	tests := []struct {
		in  []byte
		exp int
	}{
		{nil, 0},
		{[]byte{0, 0, 0, 0}, -25144152},
		{[]byte{1, 0, 0, 0}, -24390581},
		{[]byte{1, 0, 0, 1}, -24390580},
		{[]byte{1, 0, 0xff, 1}, -24367375},
		{[]byte("<*e7"), 20346417 + 74529 + 6188 + 22},
	}

	for _, test := range tests {
		got := decodeBase91(test.in)
		if got != test.exp {
			t.Errorf("decodeBase64(%v) = %d, want %d", test.in, got, test.exp)
		}
	}
}

func TestPositionParsing(t *testing.T) {
	x, err := compressedParser("123")
	if err == nil {
		t.Errorf("Expected error on three bytes compressed, got %v", x)
	}

	x, err = uncompressedParser("123")
	if err == nil {
		t.Errorf("Expected error on three bytes uncompressed, got %v", x)
	}
}

func TestInvalidPosition(t *testing.T) {
	testBodies := []string{
		"@",
		"@1234568",
		"!",
		"!1",
		";",
		";123456789012345678",
	}
	for i, testFrame := range testBodies {
		t.Run(fmt.Sprintf("InvalidPosition[%d]", i), func(t *testing.T) {
			parsedFrame := ParseFrame("SOURCE>DESTINATION,PATH:" + testFrame)
			_, err := parsedFrame.Body.Position()
			if err == nil {
				t.Fatalf("Parsing %q: expecting any error, go not error", testFrame)
			}
		})
	}
}
