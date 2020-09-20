package users

import (
	"errors"
)

// UserUseCases represents the users use cases.
type UserUseCases struct {
	db UserDBInterface
}

// NewUserUseCases returns a new user use case.
func NewUserUseCases(db UserDBInterface) UserUseCases {
	return UserUseCases{db: db}
}

// GetUser returns an user.
// The user is created if it does not exist.
//   The external service that calls our service will only pass valid users.
//   That's why we can create the user without performing any checks.
func (uc UserUseCases) GetUser(userID UserID) (entity *User, created bool, err error) {
	entity, err = uc.db.GetByID(userID)
	if err != nil {
		return nil, false, err
	}
	if entity != nil {
		return entity, false, nil
	}

	// TODO: The use case should be resposible to setup the CreatedAt and UpdatedAt,
	// but at this time we are leaving this job for the ORM because of problems with
	// tests using current time.
	newEnt := User{ID: userID}
	entity, err = uc.db.Insert(newEnt)
	if err == nil {
		return entity, false, nil
	}

	if errors.Is(err, ErrUserAlreadyExists) {
		// The user was inserted between the GetByID and Insert.
		// We will try to get this user again.
		entity, err = uc.db.GetByID(userID)
		if err != nil {
			return nil, false, err
		}
		if entity != nil {
			return entity, false, nil
		}
	}
	// Rare possibility: The user does not exist and we try to get it.
	// At the same time it was inserted and we get a "already exists error".
	// So we try to get this user again, but at this time it was deleted.
	return nil, false, err
}
