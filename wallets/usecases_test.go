package wallets_test

import (
	"fmt"
	"home-broker/money"
	userstests "home-broker/tests/users"
	userstestsmocks "home-broker/tests/users/mocks"
	walletstests "home-broker/tests/wallets"
	"home-broker/tests/wallets/mocks"
	"home-broker/users"
	"home-broker/wallets"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetWallet_WalletExists_ReturnsWallet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedEnt := walletstests.GetWalletWithDeletedAt()

	mockUserDB := userstestsmocks.NewMockUserDBInterface(mockCtrl)
	userUC := users.NewUserUseCases(mockUserDB)

	mockDB := mocks.NewMockWalletDBInterface(mockCtrl)

	retEnt := expectedEnt
	mockDB.EXPECT().
		GetByUserID(expectedEnt.UserID).
		Return(&retEnt, nil)

	uc := wallets.NewWalletUseCases(mockDB, userUC)

	entity, created, userCreated, err := uc.GetWallet(expectedEnt.UserID)
	if err != nil {
		t.Fatal(err)
	}

	err = walletstests.CheckWallets(*entity, expectedEnt)
	if err != nil {
		t.Error(err)
	}

	if created {
		t.Errorf("the wallet was created, expected as not created")
	}
	if userCreated {
		t.Errorf("the user was created, expected as not created")
	}
}

func TestGetWallet_WalletDoesNotExists_WalletCreated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedEnt := walletstests.GetWallet()

	mockUserDB := userstestsmocks.NewMockUserDBInterface(mockCtrl)
	userUC := users.NewUserUseCases(mockUserDB)

	mockDB := mocks.NewMockWalletDBInterface(mockCtrl)

	expectedUserEnt := userstests.GetEntity()

	mockDB.EXPECT().
		GetByUserID(expectedEnt.UserID).
		Return(nil, nil)

	retUserEnt := expectedUserEnt
	mockUserDB.EXPECT().
		GetByID(expectedUserEnt.ID).
		Return(&retUserEnt, nil)

	preInsertEnt := wallets.Wallet{
		UserID:  expectedEnt.UserID,
		Balance: money.NewMoneyZero(),
	}

	retEnt := expectedEnt
	mockDB.EXPECT().
		Insert(preInsertEnt).
		Return(&retEnt, nil)

	uc := wallets.NewWalletUseCases(mockDB, userUC)

	entity, created, _, err := uc.GetWallet(expectedEnt.UserID)
	if err != nil {
		t.Error(err)
	}

	err = walletstests.CheckWallets(*entity, expectedEnt)
	if err != nil {
		t.Error(err)
	}
	if !created {
		t.Errorf("the wallet not created, expected as created")
	}
}

func TestGetWallet_BetweenGetByUserIDAndInsertCalls_ReturnsWallet(t *testing.T) {
	// This test the event of a wallet inserted between the GetByUserID and Insert calls.
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedEnt := walletstests.GetWallet()

	mockUserDB := userstestsmocks.NewMockUserDBInterface(mockCtrl)
	userUC := users.NewUserUseCases(mockUserDB)

	mockDB := mocks.NewMockWalletDBInterface(mockCtrl)

	expectedUserEnt := userstests.GetEntity()

	mockDB.EXPECT().
		GetByUserID(expectedEnt.UserID).
		Return(nil, nil)

	returnedUseEnt := expectedUserEnt
	mockUserDB.EXPECT().
		GetByID(expectedUserEnt.ID).
		Return(&returnedUseEnt, nil)

	preInsertEnt := wallets.Wallet{
		UserID:  expectedEnt.UserID,
		Balance: money.NewMoneyZero(),
	}
	mockDB.EXPECT().
		Insert(preInsertEnt).
		Return(nil, fmt.Errorf("%w: User ID %d", wallets.ErrWalletAlreadyExists, expectedEnt.UserID))

	returnedEnt := expectedEnt
	mockDB.EXPECT().
		GetByUserID(expectedEnt.UserID).
		Return(&returnedEnt, nil)

	uc := wallets.NewWalletUseCases(mockDB, userUC)

	entity, created, _, err := uc.GetWallet(expectedEnt.UserID)
	if err != nil {
		t.Error(err)
	}

	err = walletstests.CheckWallets(*entity, expectedEnt)
	if err != nil {
		t.Error(err)
	}
	if created {
		t.Errorf("the wallet created, expected as not created")
	}
}

func TestAddFunds(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserDB := userstestsmocks.NewMockUserDBInterface(mockCtrl)
	userUC := users.NewUserUseCases(mockUserDB)

	mockDB := mocks.NewMockWalletDBInterface(mockCtrl)

	expectedEnt := walletstests.GetWallet()

	incValue, err := money.NewMoneyFromFloatString("789.789")
	if err != nil {
		t.Error(err)
	}
	balance, err := money.NewMoneyFromFloatString("333.333333")
	if err != nil {
		t.Error(err)
	}

	expectedEnt.Balance = balance + incValue

	returnedEnt := expectedEnt

	retEnt := expectedEnt

	mockDB.EXPECT().
		GetByUserID(expectedEnt.UserID).
		Return(&retEnt, nil)

	mockDB.EXPECT().
		IncBalanceByUserID(expectedEnt.UserID, incValue).
		Return(&returnedEnt, nil)

	uc := wallets.NewWalletUseCases(mockDB, userUC)

	entity, err := uc.AddFunds(expectedEnt.UserID, incValue)
	if err != nil {
		t.Error(err)
	}

	err = walletstests.CheckWallets(*entity, expectedEnt)
	if err != nil {
		t.Error(err)
	}
}

func TestAddFunds_InvalidAmount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserDB := userstestsmocks.NewMockUserDBInterface(mockCtrl)
	userUC := users.NewUserUseCases(mockUserDB)

	mockDB := mocks.NewMockWalletDBInterface(mockCtrl)

	userID := users.UserID(999)

	t.Run("ZeroAmount", func(t *testing.T) {
		incValue, err := money.NewMoneyFromFloatString("0.0")
		if err != nil {
			t.Error(err)
		}

		uc := wallets.NewWalletUseCases(mockDB, userUC)

		entity, err := uc.AddFunds(userID, incValue)
		if err == nil {
			t.Errorf("an error was expected to happen")
		}
		if entity != nil {
			t.Errorf("received entity %v, expected nil", entity)
		}
	})

	t.Run("NegativeAmount", func(t *testing.T) {
		incValue, err := money.NewMoneyFromFloatString("-789.789")
		if err != nil {
			t.Error(err)
		}

		uc := wallets.NewWalletUseCases(mockDB, userUC)

		entity, err := uc.AddFunds(userID, incValue)
		if err == nil {
			t.Errorf("an error was expected to happen")
		}
		if entity != nil {
			t.Errorf("received entity %v, expected nil", entity)
		}
	})

}
