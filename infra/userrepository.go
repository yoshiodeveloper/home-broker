package infra

import (
	"broker-dealer/domain"
)

// UserRepository is an interface that handles database commands for User entity.
type UserRepository interface {
	Repository

	// GetByID must get a user from the database by an ID.
	GetByID(id domain.UserID) (user domain.User, found bool, err error)
}
