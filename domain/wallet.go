// Package domain holds all the entities and operations on them.
package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

// Currency is a data type for monetary values.
type Currency struct {
	Number decimal.Decimal
}

// NewCurrencyFromString creates a new Currency from a string.
// The value must be a non-localized float number.
//   Good values: "9", "9.99", "-0.99", "-99.9999".
//   Bad values: "1,99", "9,999,999.99", "6.6062e+20".
func NewCurrencyFromString(value string) Currency {
	num, err := decimal.NewFromString(value)
	if err != nil {
		panic(err)
	}
	return Currency{Number: num}
}

// IsZero returns true if the number is a perfect zero value (0.0000...).
func (c *Currency) IsZero() bool {
	return c.Number.IsZero()
}

// WalletID represents the Wallet ID type.
//   This eases a future DB change.
type WalletID int64

// Wallet represents an User entity wallet.
type Wallet struct {
	ID        WalletID
	UserID    UserID
	Balance   Currency
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
