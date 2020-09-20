package wallets

import (
	"errors"
	"home-broker/money"
	"home-broker/users"
)

var (
	// ErrInvalidFundsAmount happens when the amount of a add funds is invalid.
	ErrInvalidFundsAmount = errors.New("invalid funds amount")
)

// WalletUseCases represents the wallets use cases.
type WalletUseCases struct {
	db     WalletDBInterface
	userUC users.UserUseCases
}

// NewWalletUseCases returns a new WalletUseCases.
func NewWalletUseCases(db WalletDBInterface, userUC users.UserUseCases) WalletUseCases {
	return WalletUseCases{db: db, userUC: userUC}
}

// GetWallet returns a wallet by an user ID.
// A new empty wallet is created if the wallet does not exist.
// In this case the user is also created because a missing wallet can be a missing user as well.
//   The external service that calls our service will only pass valid users.
//   We are considering that this external service is reliable and we do not need to verify the user.
func (uc WalletUseCases) GetWallet(userID users.UserID) (entity *Wallet, created bool, userCreated bool, err error) {
	entity, err = uc.db.GetByUserID(userID)
	if err != nil {
		return nil, false, false, err
	}
	if entity != nil {
		// wallet found by user id
		return entity, false, false, nil
	}

	// At this point the wallet was not found.
	// Maybe the user does not exist as well.

	// uc.userUC.GetUser will create the user if needed.
	_, userCreated, err = uc.userUC.GetUser(userID)
	if err != nil {
		return nil, false, false, err
	}

	// TODO: The use case should be resposible to setup the CreatedAt and UpdatedAt,
	// but at this time we are leaving this job for the ORM because of problems with
	// tests using current time.

	balance := money.NewMoneyZero()

	newWallet := Wallet{
		UserID:  userID,
		Balance: balance,
	}

	entity, err = uc.db.Insert(newWallet)
	if err == nil {
		return entity, true, userCreated, nil
	}

	if errors.Is(err, ErrWalletAlreadyExists) {
		// Maybe the wallet was inserted by other process in the meantime.
		entity, err = uc.db.GetByUserID(userID)
		if err == nil {
			return entity, false, userCreated, nil
		}
	}

	// Unknown error.
	return nil, false, userCreated, err
}

// AddFunds increments the wallet funds (balance).
// The wallet will be created if it does not exist.
func (uc WalletUseCases) AddFunds(userID users.UserID, amount money.Money) (entity *Wallet, err error) {
	if amount <= 0 {
		return nil, ErrInvalidFundsAmount
	}

	_, _, _, err = uc.GetWallet(userID) // forces wallet/user creation
	if err != nil {
		return nil, err
	}

	// We leave this job to the ORM as it can optimize this process in a single call.
	entity, err = uc.db.IncBalanceByUserID(userID, amount)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
