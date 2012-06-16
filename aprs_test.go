package aprs

import (
	"encoding/json"
	"math"
	"os"
	"testing"
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

func assertEpsilon(t *testing.T, field string, expected, got float64) {
	if math.Abs(got-expected) > 0.001 {
		t.Fatalf("Expected %v for %v, got %v -- of by %v",
			expected, field, got, math.Abs(got-expected))
	}
}

func TestAPRS(t *testing.T) {
	v := ParseAPRSMessage(CHRISTMAS_MSG)
	assert(t, "Source", v.Source, "KG6HWF")
	assert(t, "Dest", v.Dest, "APX200")
	assert(t, "len(Path)", len(v.Path), 2)
	assert(t, "Path[0]", v.Path[0], "WIDE1-1")
	assert(t, "Path[1]", v.Path[1], "WIDE2-1")
	assert(t, "Body", string(v.Body), "=3722.1 N/12159.1 W-Merry Christmas!")

	lat, lon, err := v.Body.Position()
	if err != nil {
		t.Fatalf("Couldn't parse body position:  %v", err)
	}

	assertEpsilon(t, "lat", 37.368333333333, lat)
	assertEpsilon(t, "lon", -121.985, lon)

	assert(t, "String()", v.String(), CHRISTMAS_MSG)
}

type SampleDoc struct {
	Src           string                 `json:"src"`
	Result        map[string]interface{} `Json:"result"`
	Failed        int                    `json:"failed"`
	Misunderstood bool
}

func assertLatLon(t *testing.T, lat, lon float64, doc SampleDoc) {
	slat, haslat := doc.Result["latitude"].(float64)
	slon, haslon := doc.Result["longitude"].(float64)
	if !(haslat && haslon) {
		return
	}
	if math.Abs(lat-slat) > 0.001 || math.Abs(lon-slon) > 0.001 {
		t.Fatalf("Error parsing lat/lon from %v, got %v,%v; expected %v,%v",
			doc.Src, lat, lon, slat, slon)
	}

}

func negAssertLatLon(t *testing.T, lat, lon float64, doc SampleDoc) {
	slat, haslat := doc.Result["latitude"].(float64)
	slon, haslon := doc.Result["longitude"].(float64)
	if !(haslat && haslon) {
		return
	}
	if !(math.Abs(lat-slat) > 0.001 || math.Abs(lon-slon) > 0.001) {
		t.Fatalf("Expected to fail parsing lat/lon from %v, got %v,%v; expected %v,%v",
			doc.Src, lat, lon, slat, slon)
	}
}

func TestFAP(t *testing.T) {
	minSuccess := 18

	var samples []SampleDoc
	r, err := os.Open("sample.json")
	if err != nil {
		t.Fatalf("Error opening sample.json")
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&samples)
	if err != nil {
		t.Fatalf("Error reading JSON: %", err)
	}
	t.Logf("Found %d messages", len(samples))

	positions := 0
	misunderstood := 0

	for _, sample := range samples {
		if sample.Failed != 1 {
			v := ParseAPRSMessage(sample.Src)
			assert(t, "Source", v.Source, sample.Result["srccallsign"])
			assert(t, "Dest", v.Dest, sample.Result["dstcallsign"])
			assert(t, "Body", string(v.Body), sample.Result["body"])

			if sample.Misunderstood {
				misunderstood++
				lat, lon, err := v.Body.Position()
				if err == nil {
					negAssertLatLon(t, lat, lon, sample)
				}
			} else {
				lat, lon, err := v.Body.Position()
				if err == nil {
					assertLatLon(t, lat, lon, sample)
					positions++
				}
			}

		}
	}

	if positions < minSuccess {
		t.Fatalf("Expected to pass at least %v position tests, got %v",
			minSuccess, positions)
	}

	t.Logf("Found %v positions", positions)
	t.Logf("Misunderstood %v", misunderstood)
}

func TestDecodeBase91(t *testing.T) {
	v := decodeBase91([]byte("<*e7"))
	expected := 20346417 + 74529 + 6188 + 22
	if v != expected {
		t.Fatalf("Expected %v, got %v", expected, v)
	}
}
