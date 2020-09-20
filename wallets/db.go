package wallets

import (
	"errors"
	"home-broker/money"
	"home-broker/users"
)

var (
	// ErrWalletDoesNotExist happens when the wallet record does not exist.
	ErrWalletDoesNotExist = errors.New("wallet does not exist")

	// ErrWalletAlreadyExists happens when wallet record already exists.
	ErrWalletAlreadyExists = errors.New("wallet already exists")
)

// WalletDBInterface is an interface that handles database commands for Wallet entity.
type WalletDBInterface interface {
	// GetByUserID must return a wallet by an user ID.
	// A nil entity will be returned if it does not exist.
	GetByUserID(userID users.UserID) (*Wallet, error)

	// Insert must insert a new wallet.
	// A nil entity will be returned if an error occurs.
	// The following errors can happen: ErrWalletAlreadyExists, ErrUserDoesNotExist.
	Insert(entity Wallet) (*Wallet, error)

	// IncBalanceByUserID must increment or decrement the balance field by a user ID.
	// An updated entity entity will be returned or nil if it does not exist.
	IncBalanceByUserID(userID users.UserID, amount money.Money) (*Wallet, error)
}
