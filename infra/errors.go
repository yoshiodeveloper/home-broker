package infra

import "errors"

var (
	// ErrUserDoesNotExist user record does not exist error.
	ErrUserDoesNotExist = errors.New("user does not exist")

	// ErrWalletDoesNotExist wallet record does not exist error.
	ErrWalletDoesNotExist = errors.New("wallet does not exist")

	// ErrUserAlreadyExists user record already exists error.
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrWalletAlreadyExists wallet record already exists error.
	ErrWalletAlreadyExists = errors.New("wallet already exists")
)
