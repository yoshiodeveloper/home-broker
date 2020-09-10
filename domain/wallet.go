// Package domain holds all the entities and operations on them.
package domain

import (
	"fmt"
	"math"
)

const (
	// CurrencyDecimalPrecision is the maximum decimal precision of Currency.
	// If you need change the precision just need change this value.
	CurrencyDecimalPrecision int64 = 4
)

// maxCurrencyFraction is the highest value for a monetary fraction.
// If you need change the Currency precision, you just change CurrencyDecimalPrecision.
var maxCurrencyFraction int64 = int64(math.Pow10(int(CurrencyDecimalPrecision)))

// Currency is a data type for monetary values.
// It stores the currency value as an int64 where the firsts 4 rightest digits are the decimal value.
// For example, the value 12.3456 must be stored as 123456.
type Currency int64

// NewCurrency creates a new Currency value from an integer and a decimal part.
// Notice that the decimal part is not the "cents". For example, 99 is ".0099" and not ".99" cents.
// In this case, to express a ".99" cents you need set decimalPart as 9900.
// This is necessary to represent fractions like ".0001" ($99.0001 is not the same as $99.01).
// If you need to express a value as -0.99 you must use a negative decimalPart, because it is impossible to integerPart be a "negative zero".
func NewCurrency(integerPart int64, decimalPart int64) Currency {
	isNegative := false
	if integerPart < 0 {
		isNegative = true
		// integerPart must be a positive value to be added with decimalPart.
		// Ex.:
		//   90000 + 1000 = 91000 (or $9.1000)
		//   but
		//   -90000 + 1000 = -89000 (this is not -91000 or -$9.1000)
		integerPart = -integerPart
	}
	if decimalPart < 0 {
		isNegative = true
		// It cannot be negative or an infinite loop will occur.
		decimalPart = -decimalPart
	}

	// Truncates the decimal part. Ex., ".1234567" becomes ".1234".
	for decimalPart > maxCurrencyFraction {
		decimalPart /= 10
	}

	c := Currency(integerPart*maxCurrencyFraction + decimalPart)
	if isNegative {
		c = -c
	}

	return c
}

// AsFloat64 returns the currency value as a float64.
// Do not use this to do calculations.
func (c Currency) AsFloat64() float64 {
	return float64(c) / float64(maxCurrencyFraction)
}

// String returns Currency as string.
func (c Currency) String() string {
	isNegative := false
	if c < 0 {
		isNegative = true
		// To positive. It avoids problems with mod.
		c = -c
	}
	integerPart := int64(c) / int64(maxCurrencyFraction)
	fractionPart := int64(c) % int64(maxCurrencyFraction)
	fmtStr := "%d.%0*d"
	if isNegative {
		fmtStr = "-%d.%0*d"
	}
	return fmt.Sprintf(fmtStr, integerPart, CurrencyDecimalPrecision, fractionPart)
}
