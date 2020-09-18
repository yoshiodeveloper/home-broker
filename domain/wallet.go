// Package domain holds all the entities and operations on them.
package domain

import (
	"time"
)

// Wallet represents an User entity wallet.
type Wallet struct {
	ID        WalletID
	UserID    UserID
	Balance   Money
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
