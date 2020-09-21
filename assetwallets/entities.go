package assetwallets

import (
	"home-broker/assets"
	"home-broker/users"
	"time"
)

// AssetWalletID represents the asset wallet ID type.
type AssetWalletID int64

// AssetWallet represents an user entity wallet for assets.
type AssetWallet struct {
	ID        AssetWalletID    `json:"id"`
	UserID    users.UserID     `json:"user_id"`
	AssetID   assets.AssetID   `json:"asset_id"`
	Balance   assets.AssetUnit `json:"balance"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt time.Time        `json:"-"`
}
