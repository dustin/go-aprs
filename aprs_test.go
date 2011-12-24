package aprs

import (
	"encoding/json";
	"os";
	"testing";
)

const CHRISTMAS_MSG string = "KG6HWF>APX200,WIDE1-1,WIDE2-1:=3722.1 N/12159.1 W-Merry Christmas!"

type aprsTest struct {
	in string
}

func assert(t *testing.T, name string, got interface{}, expected interface{}) {
	if got != expected {
		t.Fatalf("Expected %s for %s, got %s", expected, name, got)
	}
	// t.Logf("Looks like %s was %s", name, expected)
}

func TestAPRS(t *testing.T) {
	v := ParseAPRSMessage(CHRISTMAS_MSG)
	assert(t, "Source", v.Source, "KG6HWF")
	assert(t, "Dest", v.Dest, "APX200")
	assert(t, "len(Path)", len(v.Path), 2)
	assert(t, "Path[0]", v.Path[0], "WIDE1-1")
	assert(t, "Path[1]", v.Path[1], "WIDE2-1")
	assert(t, "Body", v.Body, "=3722.1 N/12159.1 W-Merry Christmas!")

	assert(t, "ToString()", v.ToString(), CHRISTMAS_MSG)
}

type SampleDoc struct {
	Src string `json:"src"`
	Result map[string]interface{} `Json:"result"`
	Failed int `json:"failed"`
}

func TestFAP(t *testing.T) {
	var samples []SampleDoc
	r, oerr := os.Open("sample.json")
	if oerr != nil {
		t.Fatalf("Error opening sample.json")
	}
	jd := json.NewDecoder(r)
	jd.Decode(&samples)
	t.Logf("Found %d messages", len(samples))

	for _, sample := range(samples) {
		if sample.Failed != 1 {
			v := ParseAPRSMessage(sample.Src)
			assert(t, "Source", v.Source, sample.Result["srccallsign"])
			assert(t, "Dest", v.Dest, sample.Result["dstcallsign"])
			assert(t, "Body", v.Body, sample.Result["body"])
		}
	}
}
