package aprs

import "testing"

func TestSymbol(t *testing.T) {
	tests := []struct {
		in Symbol
		s  string
	}{
		{Symbol{'/', '\''}, `{/': Plane sm}`},
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
		{[]byte{0, 0, 0, 0}, -25120856},
		{[]byte{1, 0, 0, 0}, -24367285},
		{[]byte{1, 0, 0, 1}, -24367284},
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
