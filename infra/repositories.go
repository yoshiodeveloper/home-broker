package infra

import "errors"

var (
	// ErrUserDoesNotExist happens when the user record does not exist.
	ErrUserDoesNotExist = errors.New("user does not exist")

	// ErrWalletDoesNotExist happens when the wallet record does not exist.
	ErrWalletDoesNotExist = errors.New("wallet does not exist")

	// ErrUserAlreadyExists happens when the user record already exists.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrWalletAlreadyExists happens when wallet record already exists.
	ErrWalletAlreadyExists = errors.New("wallet already exists")
)

// Repositories holds all repositories.
// Repositories know how to interact with the database.
// They are passed to the use cases.
type Repositories struct {
	UserRepo   UserRepoI
	WalletRepo WalletRepoI
}

// RepositoryI is an interface for a repository.
// Must be used by each repository interface.
type RepositoryI interface{}
