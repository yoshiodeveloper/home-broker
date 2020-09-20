package assets

import "time"

type (
	// AssetID represents the Asset ID type (or "stock ticker").
	AssetID string

	// ExchangeID represents the stock exchange ID (or broker).
	ExchangeID string
)

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
