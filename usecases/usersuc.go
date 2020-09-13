package usecases

import (
	"home-broker/domain"
	"home-broker/infra"
)

// UserUC represents the users use cases.
type UserUC struct {
	UseCase
	repos *infra.Repositories
}

// NewUserUC returns a new WalletUseCase.
func NewUserUC(repos *infra.Repositories) *UserUC {
	return &UserUC{repos: repos}
}

// GetUser returns an user.
// The user is created if it does not exist.
//   The external service that calls our service will only pass valid users.
//   That's why we can create the user without performing any checks.
func (uc *UserUC) GetUser(userID domain.UserID) (user domain.User, created bool) {
	user, found := uc.repos.UserRepo.GetByID(userID)
	if !found {
		user = domain.User{ID: userID}
		user = uc.repos.UserRepo.Insert(user)
		created = true
	}
	return user, created
}
