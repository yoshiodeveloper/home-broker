package tests

import (
	"fmt"
	"home-broker/money"
	userstests "home-broker/tests/users"
	"home-broker/wallets"
	"time"
)

var (
	// BaseTime a base time for wallets.
	BaseTime = time.Date(2020, time.Month(2), 0, 0, 0, 0, 0, time.UTC)
)

// GetWallet returns a wallet entity.
func GetWallet() wallets.Wallet {
	user := userstests.GetEntity()
	balance := money.Money(999999999999999)
	entity := wallets.Wallet{
		ID:        wallets.WalletID(9999),
		UserID:    user.ID,
		Balance:   balance,
		CreatedAt: BaseTime,
		UpdatedAt: BaseTime.Add(time.Hour * 2),
		DeletedAt: time.Time{},
	}
	return entity
}

// GetWalletWithDeletedAt returns a wallet entity with DeletedAt set.
func GetWalletWithDeletedAt() wallets.Wallet {
	entity := GetWallet()
	entity.DeletedAt = BaseTime.Add(time.Hour * 3)
	return entity
}

// CheckWallets compares if two wallets are equals.
func CheckWallets(a wallets.Wallet, b wallets.Wallet) error {
	if a.ID != b.ID {
		return fmt.Errorf("wallet.ID is %v, expected %v", a.ID, b.ID)
	}
	if a.UserID != b.UserID {
		return fmt.Errorf("wallet.UserID is %v, expected %v", a.UserID, b.UserID)
	}
	if a.Balance != b.Balance {
		return fmt.Errorf("wallet.Balance is %v, expected %v", a.Balance, b.Balance)
	}
	if a.CreatedAt != b.CreatedAt {
		return fmt.Errorf("wallet.CreatedAt is %v, expected %v", a.CreatedAt, b.CreatedAt)
	}
	if a.UpdatedAt != b.UpdatedAt {
		return fmt.Errorf("wallet.UpdatedAt is %v, expected %v", a.UpdatedAt, b.UpdatedAt)
	}
	if a.DeletedAt != b.DeletedAt {
		return fmt.Errorf("wallet.DeletedAt is %v, expected %v", a.DeletedAt, b.DeletedAt)
	}
	return nil
}
