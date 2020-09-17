// Package domain holds all the entities and operations on them.
package domain

import (
	"math"
	"time"

	"github.com/shopspring/decimal"
)

// MoneyDecimalPlaces indicates the precision of a Money type.
// This means thenumbers before the point separator.
const MoneyDecimalPlaces int = 6

// Money is a data type for monetary values.
// There is decimal part. Ex:
// $10.55 -> Money(10550000)
// $999.99 -> Money(999990000)
// $1.999999 -> Money(1999999)
type Money int64

// NewMoneyFromFloatString creates a new Money from a string.
// The value must be a non-localized float number.
//   Good values: "9", "9.99", "-0.99", "-99.9999".
//   Bad values: "1,99", "9,999,999.99", "6.6062e+20".
func NewMoneyFromFloatString(value string) Money {
	num, err := decimal.NewFromString(value)
	if err != nil {
		panic(err)
	}
	mult := math.Pow10(MoneyDecimalPlaces)
	num = num.Mul(decimal.NewFromFloat(mult))
	m := Money(num.IntPart())
	return m
}

/*
// Money is a data type for monetary values.
type Money struct {
	Number decimal.Decimal
}

// NewMoneyFromString creates a new Money from a string.
// The value must be a non-localized float number.
//   Good values: "9", "9.99", "-0.99", "-99.9999".
//   Bad values: "1,99", "9,999,999.99", "6.6062e+20".
func NewMoneyFromString(value string) Money {
	num, err := decimal.NewFromString(value)
	if err != nil {
		panic(err)
	}
	m := Money{Number: num}
	m.truncate()
	return m
}

// truncate truncates the decimal value to the maximum decimal precision allowed.
func (c *Money) truncate() {
	c.Number = c.Number.Truncate(MoneyDecimalPlaces)
}

// String converts to string.
func (c *Money) String() string {
	return c.Number.StringFixedBank(MoneyDecimalPlaces)
}

// NewMoneyFromFloat64 creates a new Money from a float64.
// Be aware the this can cause data loss.
func NewMoneyFromFloat64(value float64) Money {
	number := decimal.NewFromFloat(value)
	m := Money{Number: number}
	m.truncate()
	return m
}

// AsFloat64 returns as float64.
// Be aware the this can cause data loss.
func (c *Money) AsFloat64() float64 {
	v, _ := c.Number.Float64()
	return v
}

// IsZero returns true if the number is a perfect zero value (0.0000...).
func (c *Money) IsZero() bool {
	return c.Number.IsZero()
}

// Equal compares with another Money.
func (c *Money) Equal(v Money) bool {
	return c.Number.Equal(v.Number)
}

// LessThan returns true when this is less than v.
func (c *Money) LessThan(v Money) bool {
	return c.Number.LessThan(v.Number)
}

// GreaterThan returns true when this is greater than v.
func (c *Money) GreaterThan(v Money) bool {
	return c.Number.GreaterThan(v.Number)
}

// Add increments (or decrements) the value.
func (c *Money) Add(v Money) Money {
	return Money{Number: c.Number.Add(v.Number)}
}
*/

// WalletID represents the Wallet ID type.
//   This eases a future DB change.
type WalletID int64

// Wallet represents an User entity wallet.
type Wallet struct {
	ID        WalletID
	UserID    UserID
	Balance   Money
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
