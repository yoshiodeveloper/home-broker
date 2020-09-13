package infra

import "home-broker/domain"

// Repositories holds all repositories.
// Repositories know how to interact with the database.
type Repositories struct {
	UserRepo   UserRepo
	WalletRepo WalletRepo
}

// Repository is an interface for a repository.
type Repository interface {
}

// UserRepo is an interface that handles database commands for User entity.
type UserRepo interface {
	Repository

	// GetByUserID must return an user by ID.
	GetByID(id domain.UserID) (user domain.User, found bool)

	// Insert inserts a new user.
	// The repository must handle when a record already exists.
	Insert(user domain.User) (newUser domain.User)
}

// WalletRepo is an interface that handles database commands for Wallet entity.
type WalletRepo interface {
	Repository

	// GetByUserID must return a wallet by an user ID.
	GetByUserID(id domain.UserID) (wallet domain.Wallet, found bool)

	// Insert inserts a new wallet.
	// The repository must handle when a record already exists.
	Insert(wallet domain.Wallet) (newWallet domain.Wallet)
}
