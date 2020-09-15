package usecases

import (
	"errors"
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
func (uc *WalletUC) GetWallet(userID domain.UserID) (entity *domain.Wallet, created bool, userCreated bool, err error) {
	repo := uc.repos.WalletRepo
	entity, err = repo.GetByUserID(userID)
	if err != nil {
		return nil, false, false, err
	}
	if entity != nil {
		// wallet found by user id
		return entity, false, false, nil
	}

	// Wallet not found.
	// Maybe user does not exist as well.

	// UserUC.GetUser will create the user if needed.
	userUC := NewUserUC(uc.repos)
	_, userCreated, err = userUC.GetUser(userID)
	if err != nil {
		return nil, false, false, err
	}

	// TODO: The use case should be resposible to setup the CreatedAt and UpdatedAt,
	// but at this time we are leaving this job for the ORM because of problems with
	// tests using current time.

	newWallet := domain.Wallet{
		UserID:  userID,
		Balance: domain.NewMoneyFromString("0.0"),
	}
	entity, err = repo.Insert(newWallet)
	if err == nil {
		return entity, true, userCreated, nil
	}

	if errors.Is(err, infra.ErrWalletAlreadyExists) {
		// Maybe the wallet was inserted by other process in the meantime.
		entity, err = repo.GetByUserID(userID)
		if err == nil {
			return entity, false, userCreated, nil
		}
	}
	// Unknown error.
	return nil, false, userCreated, err
}

// IncBalance increments or decrements the wallet funds (balance).
// The retuned entity can be nil if the user or wallet does not exist.
func (uc *WalletUC) IncBalance(userID domain.UserID, amount domain.Money) (entity *domain.Wallet, err error) {
	repo := uc.repos.WalletRepo
	entity, err = repo.IncBalanceByUserID(userID, amount)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
