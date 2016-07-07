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
