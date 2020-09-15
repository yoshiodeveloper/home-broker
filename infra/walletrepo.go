package infra

import "home-broker/domain"

// WalletRepoI is an interface that handles database commands for Wallet entity.
type WalletRepoI interface {
	RepositoryI

	// GetByUserID must return a wallet by an user ID.
	// A nil entity will be returned if it does not exist.
	GetByUserID(userID domain.UserID) (*domain.Wallet, error)

	// Insert must insert a new wallet.
	// A nil entity will be returned if an error occurs.
	// The following errors can happen: ErrWalletAlreadyExists, ErrUserDoesNotExist.
	Insert(entity domain.Wallet) (*domain.Wallet, error)

	// IncBalanceByUserID must increment or decrement the balance field by a user ID.
	// An updated entity entity will be returned or nil if it does not exist.
	IncBalanceByUserID(userID domain.UserID, amount domain.Money) (*domain.Wallet, error)
}
