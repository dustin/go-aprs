package main

import (
	"math"
	"testing"
)

func TestPointDistance(t *testing.T) {
	p1 := Point{40.6892, -74.0444}
	p2 := Point{48.8583, 2.2945}

	d := p1.Distance(p2)
	exp := 5837.42231774235
	if math.Abs(d-exp) > .001 {
		t.Fatalf("Expected %v, got %v", exp, d)
	}
}
