package domain

import "time"

// Asset represents an asset (ex: PETR4).
type Asset struct {
	ID         AssetID    // code (or "stock ticker", as "PETR4")
	Name       string     // human name (ex. "Petrobrás")
	ExchangeID ExchangeID // ex. "B3" (BM&FBOVESPA), "NASDAQ"
	// TODO: Include a "type".
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// NewAsset returns a new Asset.
func NewAsset(id AssetID, exchangeID ExchangeID) Asset {
	return Asset{ID: id, ExchangeID: exchangeID}
}
