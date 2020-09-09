package domain

import (
	"fmt"
)

const (
	// MaxCurrencyFraction is the highest value for a monetary fraction.
	// If you need change the Currency precision you need change this value. To change it you need set as "10 ^ DecimalPrecisionYouWant".
	// Currently the value of 10000 allows the use of 4 decimal places.
	// This is hard-coded to avoid pow calculations.
	MaxCurrencyFraction int64 = 10000
)

// Currency is a data type for monetary values.
// Do not use this directly. Use NewCurrency to create a Currency value.
// It stores the currency value as an int64 where the firsts 4 rightest digits are the decimal value.
// For example, the value 12.3456 is stored as 123456. The MaxCurrencyFraction is used to extract the integer and fractional parts.
type Currency int64

// NewCurrency creates a new Currency value.
// Notice that fractionPart is expressed as fraction of MaxCurrencyFraction and not as "cents".
//  For example, when MaxCurrencyFraction is 10000, the fractionPart as 99 is ".0099" and not ".99" cents.
//  In this case, to express a ".99" cents you need set fractionPart as 9900.
//  This is necessary to represent fractions like ".0001", as $46.0001 (it is not the same as $46.01).
func NewCurrency(unitPart int64, fractionPart int64) Currency {
	isNegative := false
	if unitPart < 0 {
		isNegative = true
		// unitPart must be a positive value to be added with fractionPart.
		// Ex.:
		//   90000 + 1000 = 91000 (or $9.1000)
		//   but
		//   -90000 + 1000 = -89000 (this is not -91000 or -$9.1000)
		unitPart = -unitPart
	}
	if fractionPart < 0 {
		// It cannot be negative or an infinite loop will occur.
		fractionPart = -fractionPart
	}

	// Truncates the fractional part. Ex., 1234567 becomes 1234.
	for fractionPart > MaxCurrencyFraction {
		fractionPart /= 10
	}

	c := Currency(unitPart*MaxCurrencyFraction + fractionPart)
	if isNegative {
		c = -c
	}

	return c
}

// AsFloat64 returns the currency value as a float64.
// Do not use this to do calculations.
func (c Currency) AsFloat64() float64 {
	return float64(c) / float64(MaxCurrencyFraction)
}

// String returns Currency as string.
func (c Currency) String() string {
	return fmt.Sprintf("%.4f", c.AsFloat64())
}
