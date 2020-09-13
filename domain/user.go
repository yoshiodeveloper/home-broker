package domain

import (
	"time"
)

// UserID represents the User ID type.
//   This eases a future DB change.
type UserID int64

// User represents an user.
type User struct {
	ID        UserID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
