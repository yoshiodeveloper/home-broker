package domain

import (
	"time"
)

// User represents an user.
type User struct {
	ID        UserID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
