package domain

import (
	"time"
)

// An User entity.
type User struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
