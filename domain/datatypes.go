package domain

import (
	"math"

	"github.com/shopspring/decimal"
)

// MoneyDecimalPlaces indicates the precision of a Money type.
// This means, the numbers before the point separator.
const MoneyDecimalPlaces int = 6

// AssetID represents the Asset ID type (or "stock ticker").
type AssetID string

// ExchangeID represents the stock exchange ID (or broker).
type ExchangeID string

// UserID represents the User ID type.
type UserID int64

// OrderType represents a order type.
// Use the value of OrderTypeBuy or OrderTypeSell to set this data type.
type OrderType string

// OrderID represents the Order ID type.
type OrderID int64

// ExternalOrderID represents an external order ID (ex: an order ID from a exchange).
type ExternalOrderID string

// OrderStatus represents the status of an order.
type OrderStatus int8

// WalletID represents the Wallet ID type.
type WalletID int64

// Money is a data type for monetary values.
// It is an integer to avoid float calculations issues.
// The decimal part is the rightest 6 digits.
// Ex:
//  - $10.55    ->  10550000
//  - $999.99   -> 999990000
//  - $1.999999 ->   1999999
// The precision can be changed setting up "MoneyDecimalPlaces".
type Money int64

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
