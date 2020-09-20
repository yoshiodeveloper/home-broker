package wallets

import (
	"home-broker/money"
	"home-broker/users"
	"time"
)

// WalletID represents the Wallet ID type.
type WalletID int64

// Wallet represents an User entity wallet.
type Wallet struct {
	ID        WalletID     `json:"id"`
	UserID    users.UserID `json:"user_id"`
	Balance   money.Money  `json:"balance,string"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt time.Time    `json:"-"`
}
