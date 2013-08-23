package aprs

import (
	"bufio"
	"compress/bzip2"
	"encoding/json"
	"io"
	"math"
	"os"
	"reflect"
	"testing"
)

const (
	CHRISTMAS_MSG = "KG6HWF>APX200,WIDE1-1,WIDE2-1:=3722.1 N/12159.1 W-Merry Christmas!"
)

const SAMPLE2 = `K7FED-1>APNX01,qAR,W6MSU-7:!3739.12N112132.05W#PHG5750 W1, K7FED FILL-IN LLNL S300`

var samples = []struct {
	src      string
	expected Position
}{
	{`K6LRG-C>APJI23,WIDE1-1,WIDE2-1:!3729.98ND12152.33W&RNG0060 2m Voice 145.070 +1.495 Mhz`,
		Position{37.49966666666667, -121.87216666666667, 0, Velocity{}, Symbol{'D', '&'}}},
	{`K7FED-1>APNX01,qAR,W6MSU-7:!3739.12N112132.05W#PHG5750 W1, K7FED FILL-IN LLNL S300`,
		Position{37.652, -121.534167, 0, Velocity{}, Symbol{'1', '#'}}},
	{`WINLINK>APWL2K,TCPIP*,qAC,T2LAX:;KE6AFE-10*160752z3658.  NW12202.  Wa144.910MHz 1200 R6m Public Winlink Gateway`,
		Position{36.975, -122.0416666, 2, Velocity{}, Symbol{'W', 'a'}}},
	{`KE6AFE-13>APKH2Z,TCPIP*,qAC,CORE-2:;VP@CM86XX*162000z3658.94N/12200.86W? KE6AFE-13 8800`,
		Position{36.9823333, -122.014333, 0, Velocity{}, Symbol{'/', '?'}}},
}

func assert(t *testing.T, name string, got interface{}, expected interface{}) {
	if got != expected {
		t.Fatalf("Expected %v for %v, got %v", expected, name, got)
	}
	// t.Logf("Looks like %s was %s", name, expected)
}

func assertEpsilon(t *testing.T, field string, expected, got float64) {
	if math.Abs(got-expected) > 0.0001 {
		t.Fatalf("Expected %v for %v, got %v -- off by %v",
			expected, field, got, math.Abs(got-expected))
	}
}

func TestAddressConversion(t *testing.T) {
	testaddrs := []struct {
		Src string
		Exp Address
	}{
		{"KG6HWF", Address{"KG6HWF", 0}},
		{"KG6HWF-9", Address{"KG6HWF", 9}},
	}

	for _, ta := range testaddrs {
		a := AddressFromString(ta.Src)
		if !reflect.DeepEqual(a, ta.Exp) {
			t.Fatalf("Expected %v for %v, got %v", ta.Exp, ta.Src, a)
		}

		if ta.Exp.String() != ta.Src {
			t.Fatalf("Expected to string %v as %s, got %s",
				ta.Exp, ta.Src, ta.Exp.String())
		}
	}
}

func TestCallPass(t *testing.T) {
	testaddrs := []struct {
		Call Address
		Exp  int16
	}{
		{AddressFromString("KG6HWF-9"), 22955},
		{AddressFromString("KG6HWF"), 22955},
		{AddressFromString("KE6AFE-13"), 18595},
	}

	for _, c := range testaddrs {
		if c.Call.CallPass() != c.Exp {
			t.Fatalf("Expected %v for %v, got %v",
				c.Exp, c.Call, c.Call.CallPass())
		}
	}
}

func TestAPRS(t *testing.T) {
	v := ParseAPRSData(CHRISTMAS_MSG)
	assert(t, "Source", v.Source.String(), "KG6HWF")
	assert(t, "Dest", v.Dest.String(), "APX200")
	assert(t, "len(Path)", len(v.Path), 2)
	assert(t, "Path[0]", v.Path[0].String(), "WIDE1-1")
	assert(t, "Path[1]", v.Path[1].String(), "WIDE2-1")
	assert(t, "Body", string(v.Body), "=3722.1 N/12159.1 W-Merry Christmas!")

	pos, err := v.Body.Position()
	if err != nil {
		t.Fatalf("Couldn't parse body position:  %v", err)
	}

	assertEpsilon(t, "lat", 37.3691667, pos.Lat)
	assertEpsilon(t, "lon", -121.985833, pos.Lon)
	assert(t, "ambiguity", 1, pos.Ambiguity)
	assert(t, "table", byte('/'), pos.Symbol.Table)
	assert(t, "symbol", byte('-'), pos.Symbol.Symbol)

	assert(t, "String()", v.String(), CHRISTMAS_MSG)
}

func TestInvalid(t *testing.T) {
	v := ParseAPRSData("Invalid")
	if v.IsValid() {
		t.Fatalf("Expected invalid data out of the parse")
	}
	assert(t, "Source", v.Source.String(), "")
	assert(t, "Dest", v.Dest.String(), "")
	assert(t, "len(Path)", len(v.Path), 0)
	assert(t, "Body", string(v.Body), "")

	v = ParseAPRSData("Invalid:Thing")
	if v.IsValid() {
		t.Fatalf("Expected invalid data out of the parse")
	}
	assert(t, "Source", v.Source.String(), "")
	assert(t, "Dest", v.Dest.String(), "")
	assert(t, "len(Path)", len(v.Path), 0)
	assert(t, "Body", string(v.Body), "")
}

func TestSamples(t *testing.T) {
	for _, s := range samples {
		v := ParseAPRSData(s.src)
		pos, err := v.Body.Position()
		if err != nil {
			t.Fatalf("Error getting position from %v: %v", s.src, err)
		}
		assert(t, "ambiguity", s.expected.Ambiguity, pos.Ambiguity)
		assert(t, "table", s.expected.Symbol.Table, pos.Symbol.Table)
		assert(t, "symbol", s.expected.Symbol.Symbol, pos.Symbol.Symbol)
		assert(t, "course", s.expected.Velocity.Course, pos.Velocity.Course)
		assert(t, "speed", s.expected.Velocity.Speed, pos.Velocity.Speed)
		assertEpsilon(t, "lat", s.expected.Lat, pos.Lat)
		assertEpsilon(t, "lon", s.expected.Lon, pos.Lon)
	}
}

type SampleDoc struct {
	Src           string                 `json:"src"`
	Result        map[string]interface{} `Json:"result"`
	Failed        int                    `json:"failed"`
	Misunderstood bool
}

func assertLatLon(t *testing.T, pos Position, doc SampleDoc) {
	slat, haslat := doc.Result["latitude"].(float64)
	slon, haslon := doc.Result["longitude"].(float64)
	if !(haslat && haslon) {
		return
	}
	if math.Abs(pos.Lat-slat) > 0.001 || math.Abs(pos.Lon-slon) > 0.001 {
		t.Fatalf("Error parsing lat/lon from %v, got %v; expected %v,%v",
			doc.Src, pos, slat, slon)
	}
	tbl := doc.Result["symboltable"].(string)[0]
	if pos.Symbol.Table != tbl {
		t.Fatalf("Expected symbol table %v, got %v for %v",
			tbl, pos.Symbol.Table, doc.Src)
	}
	symbol := doc.Result["symbolcode"].(string)[0]
	if pos.Symbol.Symbol != symbol {
		t.Fatalf("Expected symbol %v, got %v for %v",
			symbol, pos.Symbol.Symbol, doc.Src)
	}
	course, _ := doc.Result["course"].(float64)
	assertEpsilon(t, "course of "+doc.Src, course, pos.Velocity.Course)
	speed, _ := doc.Result["speed"].(float64)
	assertEpsilon(t, "speed of "+doc.Src, speed, pos.Velocity.Speed)
}

func negAssertLatLon(t *testing.T, pos Position, doc SampleDoc) {
	slat, haslat := doc.Result["latitude"].(float64)
	slon, haslon := doc.Result["longitude"].(float64)
	if !(haslat && haslon) {
		return
	}
	if !(math.Abs(pos.Lat-slat) > 0.001 || math.Abs(pos.Lon-slon) > 0.001) {
		t.Fatalf("Expected to fail parsing lat/lon from %v, got %v; expected %v,%v",
			doc.Src, pos, slat, slon)
	}
}

func TestFAP(t *testing.T) {
	expSuccess := 29

	var samples []SampleDoc
	r, err := os.Open("samples/faptests.json")
	if err != nil {
		t.Fatalf("Error opening sample.json")
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&samples)
	if err != nil {
		t.Fatalf("Error reading JSON: %v", err)
	}
	t.Logf("Found %d messages", len(samples))

	positions := 0
	misunderstood := 0

	for _, sample := range samples {
		if sample.Failed != 1 {
			v := ParseAPRSData(sample.Src)
			assert(t, "Source", v.Source.String(), sample.Result["srccallsign"])
			assert(t, "Dest", v.Dest.String(), sample.Result["dstcallsign"])
			assert(t, "Body", string(v.Body), sample.Result["body"])

			if sample.Misunderstood {
				misunderstood++
				pos, err := v.Body.Position()
				if err == nil {
					negAssertLatLon(t, pos, sample)
				}
				t.Logf("Misunderstood:  %s", sample.Src)
			} else {
				pos, err := v.Body.Position()
				if err == nil {
					assertLatLon(t, pos, sample)
					positions++
				}
			}

		}
	}

	if positions != expSuccess {
		t.Fatalf("Expected to pass at %v position tests, got %v",
			expSuccess, positions)
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

func getSampleLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic("Could not open sample file: " + err.Error())
	}
	defer file.Close()

	bz := bzip2.NewReader(file)
	rv := make([]string, 0, 250000)

	bio := bufio.NewReader(bz)
	bytesread := int64(0)
	done := false

	for !done {
		line, err := bio.ReadString('\n')
		switch err {
		case nil:
			rv = append(rv, line)
			bytesread += int64(len(line))
		case io.EOF:
			done = true
		default:
			panic("Could not load samples: " + err.Error())
		}
	}

	return rv
}

var largeSample []string

func loadLargeSample(b *testing.B) {
	if largeSample == nil {
		largeSample = getSampleLines("samples/large.log.bz2")
		b.ResetTimer()
	}
}

func BenchmarkMessages(b *testing.B) {
	loadLargeSample(b)
	var read int64
	for i := 0; i < b.N; i++ {
		src := largeSample[i%len(largeSample)]
		ParseAPRSData(src)
		read += int64(len(src))
	}
	b.SetBytes(read / int64(b.N))
}

func BenchmarkPositionsFromLog(b *testing.B) {
	loadLargeSample(b)
	var read int64
	for i := 0; i < b.N; i++ {
		src := largeSample[i%len(largeSample)]
		msg := ParseAPRSData(src)
		msg.Body.Position()
		read += int64(len(src))
	}
	b.SetBytes(read / int64(b.N))
}
