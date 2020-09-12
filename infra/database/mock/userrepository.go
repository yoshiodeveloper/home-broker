package mock

import (
	"broker-dealer/domain"
	"broker-dealer/infra"
)

// UserRepository handles a mocked database commands for user table.
type UserRepository struct {
	infra.UserRepository

	// GetByID gets an user from the database by an ID.
	GetByID func(id int64) (user domain.User, found bool, err error)
}

// NewUserRepository creates a new NewUserRepository.
func NewUserRepository() UserRepository {
	return UserRepository{}
}
