package aprs

import (
	"fmt"
	"github.com/dustin/go-aprs"
)

// APRS body used for crafting APRS frame
type Body struct {
	Lat    float64
	Lon    float64
	Symbol string
}

// String - Bound stringer method for type conversions
func (b *Body) String() string {
	locationString := fmt.Sprintf("%.2fN/0%.2fE", b.Lat, b.Lon)
	return fmt.Sprintf("!%s%s%s", locationString, b.Symbol, "")
}

// Info - Returns aprs.Info type for use with APRS frame body
func (b *Body) Info() aprs.Info {
	return aprs.Info(b.String())
}
