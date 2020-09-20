package users

import (
	"errors"
)

var (
	// ErrUserDoesNotExist happens when the user record does not exist.
	ErrUserDoesNotExist = errors.New("user does not exist")

	// ErrUserAlreadyExists happens when the user record already exists.
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserDBInterface must handle database commands for User entity.
type UserDBInterface interface {
	// GetByID must return an user by ID.
	// A nil entity will be returned if it does not exist.
	GetByID(id UserID) (*User, error)

	// Insert must insert a new user.
	// A nil entity will be returned if an error occurs.
	// The following errors can happen: ErrUserAlreadyExists.
	Insert(entity User) (*User, error)
}
