package tests

import (
	"fmt"
	"home-broker/users"
	"time"
)

var (
	// BaseTime is a base time for users.
	BaseTime = time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
)

// GetEntity returns an user entity.
func GetEntity() users.User {
	return users.User{
		ID:        users.UserID(999),
		CreatedAt: BaseTime,
		UpdatedAt: BaseTime.Add(time.Hour * 2),
		DeletedAt: time.Time{},
	}
}

// GetEntityWithDeletedAt returns an user entity with a DeletedAt set.
func GetEntityWithDeletedAt() users.User {
	entity := GetEntity()
	entity.DeletedAt = BaseTime.Add(time.Hour * 3)
	return entity
}

// CheckUsers compares if two users are equals.
func CheckUsers(a users.User, b users.User) error {
	if a.ID != b.ID {
		return fmt.Errorf("user.ID is %v, expected %v", a.ID, b.ID)
	}
	if !a.CreatedAt.Equal(b.CreatedAt) {
		return fmt.Errorf("user.CreatedAt is %v, expected %v", a.CreatedAt, b.CreatedAt)
	}
	if !a.UpdatedAt.Equal(b.UpdatedAt) {
		return fmt.Errorf("user.UpdatedAt is %v, expected %v", a.UpdatedAt, b.UpdatedAt)
	}
	if !a.DeletedAt.Equal(b.DeletedAt) {
		return fmt.Errorf("user.DeletedAt is %v, expected %v", a.DeletedAt, b.DeletedAt)
	}
	return nil
}
