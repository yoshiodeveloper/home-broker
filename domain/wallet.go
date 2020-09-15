// Package domain holds all the entities and operations on them.
package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// MoneyDecimalPlaces indicates the precision of a Money type.
// This means thenumbers before the point separator.
const MoneyDecimalPlaces int32 = 6

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
	return Money{Number: num}
}

// String converts to string.
func (c *Money) String() string {
	return c.Number.StringFixedBank(MoneyDecimalPlaces)
}

// NewMoneyFromFloat64 creates a new Money from a float64.
// Be aware the this can casue data loss.
func NewMoneyFromFloat64(value float64) Money {
	number := decimal.NewFromFloat(value)
	return Money{Number: number}
}

// AsFloat64 returns as float64.
// Be aware the this can casue data loss.
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

// Add increments (or decrements) the value.
func (c *Money) Add(v Money) Money {
	return Money{Number: c.Number.Add(v.Number)}
}

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
