package assetwallets

import (
	"errors"
	"home-broker/assets"
	"home-broker/users"
)

var (
	// ErrAssetWalletDoesNotExist happens when the asset wallet record does not exist.
	ErrAssetWalletDoesNotExist = errors.New("asset wallet does not exist")

	// ErrAssetWalletAlreadyExists happens when asset wallet record already exists.
	ErrAssetWalletAlreadyExists = errors.New("asset wallet already exists")
)

// AssetWalletDBInterface is an interface that handles database commands for Wallet entity.
type AssetWalletDBInterface interface {
	// GetByUserIDAssetID must return an asset wallet by an user ID and asset ID
	// A nil entity will be returned if it does not exist.
	GetByUserIDAssetID(userID users.UserID, assetID assets.AssetID) (*AssetWallet, error)

	// Insert must insert a new asset wallet.
	// A nil entity will be returned if an error occurs.
	// The following errors can happen: ErrAssetWalletAlreadyExists, ErrAssetDoesNotExist, ErrUserDoesNotExist.
	Insert(entity AssetWallet) (*AssetWallet, error)

	// IncBalanceByUserIDAssetID must increment or decrement the balance field by a user ID and asset ID.
	// An updated entity will be returned or nil if it does not exist.
	IncBalanceByUserIDAssetID(userID users.UserID, assetID assets.AssetID, amount assets.AssetUnit) (*AssetWallet, error)
}
