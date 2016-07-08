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
