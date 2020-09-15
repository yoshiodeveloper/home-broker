package infra

import "home-broker/domain"

// WalletRepoI is an interface that handles database commands for Wallet entity.
type WalletRepoI interface {
	RepositoryI

	// GetByUserID must return a wallet by an user ID.
	GetByUserID(id domain.UserID) (*domain.Wallet, error)

	// Insert must insert a new wallet.
	Insert(entity domain.Wallet) (*domain.Wallet, error)
}
