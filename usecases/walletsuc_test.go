package usecases

import (
	"home-broker/domain"
	"home-broker/infra"
	"home-broker/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

// When we ask for an user's wallet and this user does not exist, this user must be created.
//   The external service that calls our service will only pass on valid users.
//   That's why we can create that user without performing any checks.
func TestGetWallet_WalletAndUserDoNotExist_WalletAndUserAreCreated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userID := domain.UserID(99)

	userRepo := mocks.NewMockUserRepo(mockCtrl)

	preInsertedUser := domain.User{ID: userID}
	insertedUser := domain.User{ID: userID}

	userRepo.EXPECT().
		GetByID(userID).
		Return(domain.User{}, false)

	userRepo.EXPECT().
		Insert(preInsertedUser).
		Return(insertedUser)

	walletID := domain.WalletID(1234)
	balance := domain.NewCurrencyFromString("0.0")
	preInsertedWallet := domain.Wallet{
		UserID:  userID,
		Balance: balance,
	}

	insertedWallet := domain.Wallet{
		ID:      walletID,
		Balance: balance,
		UserID:  userID,
	}
	walletRepo := mocks.NewMockWalletRepo(mockCtrl)
	walletRepo.EXPECT().
		GetByUserID(userID).
		Return(domain.Wallet{}, false)

	walletRepo.EXPECT().
		Insert(preInsertedWallet).
		Return(insertedWallet)

	useCases := NewUseCases(&infra.Repositories{UserRepo: userRepo, WalletRepo: walletRepo})

	walletUC := useCases.GetWalletUC()
	wallet := walletUC.GetWallet(userID)
	if wallet.ID != walletID {
		t.Errorf("wallet ID is %v, expected %v", wallet.ID, walletID)
	}
	if wallet.UserID != userID {
		t.Errorf("wallet user ID is %v, expected %v", wallet.UserID, userID)
	}
	if !wallet.Balance.IsZero() {
		t.Errorf("wallet balance %v, expected zero", wallet.Balance)
	}
}
