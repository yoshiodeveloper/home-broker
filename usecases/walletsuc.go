package usecases

import (
	"home-broker/domain"
	"home-broker/infra"
)

// WalletUC represents the wallets use cases.
type WalletUC struct {
	repos *infra.Repositories
}

// NewWalletUC returns a new WalletUC.
func NewWalletUC(repos *infra.Repositories) *WalletUC {
	return &WalletUC{repos: repos}
}

// GetWallet returns a wallet by user ID.
// A new empty wallet is created if the wallet does not exist.
// In this case the user is also created because a missing wallet also is a missing user.
//   The external service that calls our service will only pass valid users.
//   That's why we can create that user without performing any checks.
func (uc *WalletUC) GetWallet(userID domain.UserID) domain.Wallet {
	repo := uc.repos.WalletRepo
	wallet, found := repo.GetByUserID(userID)
	if !found {
		userUC := NewUserUC(uc.repos)
		userUC.GetUser(userID)
		wallet = domain.Wallet{
			UserID:  userID,
			Balance: domain.NewCurrencyFromString("0.0"),
		}
		wallet = repo.Insert(wallet)
	}
	return wallet
}
