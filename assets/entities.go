package assets

import (
	"errors"
	"math"
	"time"

	"github.com/shopspring/decimal"
)

type (
	// AssetID represents the Asset ID type (or "stock ticker").
	AssetID string

	// ExchangeID represents the stock exchange ID (or broker).
	ExchangeID string
)

var (
	// ErrAssetDoesNotExist happens when the asset record does not exist.
	ErrAssetDoesNotExist = errors.New("asset does not exist")

	// ErrAssetAlreadyExists happens when the asset record already exists.
	ErrAssetAlreadyExists = errors.New("asset already exists")
)

// AssetUnit is a data type for assets units.
// It is an integer to avoid float calculations issues.
// The decimal part is the rightest 6 digits.
// Ex:
//  - 10.55    ->  10550000
//  - 999.99   -> 999990000
//  - 1.999999 ->   1999999
// The precision can be changed setting up "AssetUnitDecimalPlaces".
type AssetUnit int64

const (
	// AssetUnitDecimalPlaces indicates the precision of a AssetUnit type.
	// This means, the numbers before the point separator.
	AssetUnitDecimalPlaces int = 6
)

// NewAssetUnitZero returns an asset unit with zero value.
func NewAssetUnitZero() AssetUnit {
	return AssetUnit(0)
}

// NewAssetUnitFromFloatString creates a new asset unit from a string.
// The value must be a non-localized float number.
//   Good values are "9", "9.99", "-0.99", "-99.9999".
//   Bad values are "1,99", "9,999,999.99", "6.6062e+20", etc.
func NewAssetUnitFromFloatString(value string) (AssetUnit, error) {
	num, err := decimal.NewFromString(value)
	if err != nil {
		return AssetUnit(0), err
	}
	mult := math.Pow10(AssetUnitDecimalPlaces)
	num = num.Mul(decimal.NewFromFloat(mult))
	m := AssetUnit(num.IntPart())
	return m, nil
}

// Asset represents an asset (ex: PETR4).
type Asset struct {
	ID         AssetID    `json:"id"`          // code (or "stock ticker", as "PETR4")
	Name       string     `json:"name"`        // human name (ex. "Petrobrás")
	ExchangeID ExchangeID `json:"exchange_id"` // ex. "B3" (BM&FBOVESPA), "NASDAQ"
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  time.Time  `json:"-"`
	// TODO: Add a "type".
}

// NewAsset returns a new Asset.
func NewAsset(id AssetID, exchangeID ExchangeID) Asset {
	return Asset{ID: id, ExchangeID: exchangeID}
}
