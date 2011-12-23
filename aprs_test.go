package aprs

import "testing"

type aprsTest struct {
	in string
}

func assert(t *testing.T, name string, got interface{}, expected interface{}) {
	if got != expected {
		t.Fatalf("Expected %s for %s, got %s", expected, name, got)
	}
	t.Logf("Looks like %s was %s", name, expected)
}

func TestAPRS(t *testing.T) {
	v := ParseAPRSMessage("KG6HWF>APX200,WIDE1-1,WIDE2-1:=3722.1 N/12159.1 W-Merry Christmas!")
	assert(t, "Source", v.Source, "KG6HWF")
	assert(t, "Dest", v.Dest, "APX200")
	assert(t, "len(Path)", len(v.Path), 2)
	assert(t, "Path[0]", v.Path[0], "WIDE1-1")
	assert(t, "Path[1]", v.Path[1], "WIDE2-1")
	assert(t, "Comment", v.Comment, "=3722.1 N/12159.1 W-Merry Christmas!")
}
