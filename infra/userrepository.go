package infra

import (
	"broker-dealer/domain"
)

// UserRepository is an interface that handles database commands for User table.
type UserRepository interface {

	// GetDBClient must returns a pointer to the DBClient.
	GetDBClient() *DBClient

	// GetByID must get a user from the database by an ID.
	GetByID(id int64) (user domain.User, found bool, err error)
}
