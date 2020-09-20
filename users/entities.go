package users

import (
	"time"
)

// UserID represents an user ID.
type UserID int64

// User represents an user.
type User struct {
	ID        UserID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
