package infra

import "home-broker/domain"

// UserRepoI is an interface that handles database commands for User entity.
type UserRepoI interface {
	RepositoryI

	// GetByUserID must return an user by ID.
	// A nil entity will be returned if it does not exist.
	GetByID(id domain.UserID) (*domain.User, error)

	// Insert must insert a new user.
	// A nil entity will be returned if an error occurs.
	// The following errors can happen: ErrUserAlreadyExists.
	Insert(entity domain.User) (*domain.User, error)
}
