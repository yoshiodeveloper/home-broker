package assetwallets

import (
	"errors"
	"home-broker/assets"
	"home-broker/core"
	"home-broker/users"
)

var (
	// ErrInvalidFundsAmount happens when the funds amount is zero or engative.
	ErrInvalidFundsAmount = errors.New("invalid funds amount")
)

// AssetWalletUseCases represents the asset wallets use cases.
type AssetWalletUseCases struct {
	db     AssetWalletDBInterface
	userUC users.UserUseCases
}

// NewAssetWalletUseCases returns a new WalletUseCases.
func NewAssetWalletUseCases(db AssetWalletDBInterface, userUC users.UserUseCases) AssetWalletUseCases {
	return AssetWalletUseCases{db: db, userUC: userUC}
}

// GetAssetWallet returns an asset wallet by an user ID and asset ID.
// A new empty wallet is created if the wallet does not exist.
// In this case the user is also created because a missing wallet can be a missing user as well.
//   The external service that calls our service will only pass valid users.
//   We are considering that this external service is reliable and we do not need to verify the user.
func (uc AssetWalletUseCases) GetAssetWallet(userID users.UserID, assetID assets.AssetID) (entity *AssetWallet, created bool, userCreated bool, err error) {
	if userID <= 0 {
		return nil, false, false, core.NewErrValidation("Invalid user ID.")
	}
	if assetID == "" {
		return nil, false, false, core.NewErrValidation("Invalid asset ID.")
	}

	entity, err = uc.db.GetByUserIDAssetID(userID, assetID)
	if err != nil {
		return nil, false, false, err
	}
	if entity != nil {
		// wallet found by user id
		return entity, false, false, nil
	}

	// At this point the asset wallet was not found.
	// Maybe the user does not exist as well.

	// uc.userUC.GetUser will create the user if needed.
	_, userCreated, err = uc.userUC.GetUser(userID)
	if err != nil {
		return nil, false, false, err
	}

	// TODO: The use case should be resposible to setup the CreatedAt and UpdatedAt,
	// but at this time we are leaving this job for the ORM because of problems with
	// tests using current time.

	balance := assets.NewAssetUnitZero()

	newAssetWallet := AssetWallet{
		UserID:  userID,
		AssetID: assetID,
		Balance: balance,
	}

	entity, err = uc.db.Insert(newAssetWallet)
	if err == nil {
		return entity, true, userCreated, nil
	}

	switch err {
	case assets.ErrAssetDoesNotExist:
		return nil, false, userCreated, core.NewErrValidation("Asset does not exist.")
	case users.ErrUserDoesNotExist:
		return nil, false, userCreated, core.NewErrValidation("User does not exist.")
	}

	if errors.Is(err, ErrAssetWalletAlreadyExists) {
		// Maybe the wallet was inserted by other process in the meantime.
		entity, err = uc.db.GetByUserIDAssetID(userID, assetID)
		if err == nil {
			return entity, false, userCreated, nil
		}
	}

	// Unknown error.
	return nil, false, userCreated, err
}

// AddFunds increments the asset wallet funds (balance).
// The wallet will be created if it does not exist.
func (uc AssetWalletUseCases) AddFunds(userID users.UserID, assetID assets.AssetID, amount assets.AssetUnit) (entity *AssetWallet, err error) {
	if userID <= 0 {
		return nil, core.NewErrValidation("Invalid user ID.")
	}
	if assetID == "" {
		return nil, core.NewErrValidation("Invalid asset ID.")
	}
	if amount <= 0 {
		return nil, core.NewErrValidation("Invalid amount.")
	}

	_, _, _, err = uc.GetAssetWallet(userID, assetID) // forces asset wallet/user creation
	if err != nil {
		return nil, err
	}

	// We leave this job to the ORM as it can optimize this process in a single call.
	entity, err = uc.db.IncBalanceByUserIDAssetID(userID, assetID, amount)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
