package domain

import (
	"fmt"
	"math"
)

const (
	// CurrencyPrecision is the internal decimal precision for monetary values.
	CurrencyPrecision uint8 = 4
)

// Currency is a data type for monetary values.
// It stores the currency value as an int64 where the firsts 4 rightest digits (CurrencyPrecision) are the decimal value (or "cents").
// For example, the value 12.3456 is stored as 123456.
type Currency int64

// NewCurrency creates a new Currency.
func NewCurrency(integer int64, fraction int64) Currency {
	maxFraction := int64(math.Pow10(int(CurrencyPrecision)))
	for fraction > maxFraction {
		fraction = fraction / 10
	}
	return Currency(integer*maxFraction + fraction)
}

// AsFloat64 returns the currency value as a float64.
// Do not use this to do calculations.
func (c Currency) AsFloat64() float64 {
	maxFraction := float64(math.Pow10(int(CurrencyPrecision)))
	return float64(c) / maxFraction
}

// String returns Currency as string.
func (c Currency) String() string {
	return fmt.Sprintf("%.4f", c.AsFloat64())
}
