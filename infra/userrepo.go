package infra

import "home-broker/domain"

// UserRepoI is an interface that handles database commands for User entity.
type UserRepoI interface {
	RepositoryI

	// GetByUserID must return an user by ID.
	GetByID(id domain.UserID) (*domain.User, error)

	// Insert must insert a new user.
	Insert(entity domain.User) (*domain.User, error)
}
