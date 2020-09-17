package domain

import "time"

// AssetID represents the Asset ID type (or "stock ticker").
//   This eases a future DB change.
type AssetID string

// Asset represents an asset (ex: PETR4).
type Asset struct {
	ID        AssetID // code (or "stock ticker")
	Name      string  // human name
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	// TODO:
	// ExchangeID string // NASDAQ, B3(Bovespa)
	// Type string // Shares, Bonds, etc
}
