package usecases

import "home-broker/infra"

// UseCases holds all use cases.
// Use cases know the business logic.
// They also know when they should interact with services (ex. databases).
// That's why they need to receive the Repositories using dependency injection.
type UseCases struct {
	repos *infra.Repositories
}

// NewUseCases creates a new UseCases.
func NewUseCases(repos *infra.Repositories) *UseCases {
	return &UseCases{repos: repos}
}

// GetUserUC return the users use cases.
func (uc *UseCases) GetUserUC() *UserUC {
	return NewUserUC(uc.repos)
}

// GetWalletUC return the use cases for wallets.
func (uc *UseCases) GetWalletUC() *WalletUC {
	return NewWalletUC(uc.repos)
}
