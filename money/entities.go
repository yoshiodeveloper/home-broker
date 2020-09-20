package money

import (
	"math"

	"github.com/shopspring/decimal"
)

// Money is a data type for monetary values.
// It is an integer to avoid float calculations issues.
// The decimal part is the rightest 6 digits.
// Ex:
//  - $10.55    ->  10550000
//  - $999.99   -> 999990000
//  - $1.999999 ->   1999999
// The precision can be changed setting up "MoneyDecimalPlaces".
type Money int64

const (
	// MoneyDecimalPlaces indicates the precision of a Money type.
	// This means, the numbers before the point separator.
	MoneyDecimalPlaces int = 6
)

// NewMoneyZero returns a money with zero value.
func NewMoneyZero() Money {
	return Money(0)
}

// NewMoneyFromFloatString creates a new Money from a string.
// The value must be a non-localized float number.
//   Good values are "9", "9.99", "-0.99", "-99.9999".
//   Bad values are "1,99", "9,999,999.99", "6.6062e+20", etc.
func NewMoneyFromFloatString(value string) (Money, error) {
	num, err := decimal.NewFromString(value)
	if err != nil {
		return Money(0), err
	}
	mult := math.Pow10(MoneyDecimalPlaces)
	num = num.Mul(decimal.NewFromFloat(mult))
	m := Money(num.IntPart())
	return m, nil
}
