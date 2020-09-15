package usecases_test

import (
	"fmt"
	"home-broker/domain"
	"home-broker/infra"
	"home-broker/mocks"
	"home-broker/usecases"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

var (
	walletBaseTime = time.Date(2020, time.Month(2), 0, 0, 0, 0, 0, time.UTC)
)

func GetTestWalletEntity() domain.Wallet {
	user := GetTestUserEntity()
	return domain.Wallet{
		ID:        domain.WalletID(9999),
		UserID:    user.ID,
		Balance:   domain.NewMoneyFromString("999999999.999999"),
		CreatedAt: walletBaseTime,
		UpdatedAt: walletBaseTime.Add(time.Hour * 2),
		DeletedAt: time.Time{},
	}
}

func GetTestWalletEntityWithDeletedAt() domain.Wallet {
	entity := GetTestWalletEntity()
	entity.DeletedAt = walletBaseTime.Add(time.Hour * 3)
	return entity
}

func CheckWalletsAreEquals(t *testing.T, entityA domain.Wallet, entityB domain.Wallet) {
	if entityA.ID != entityB.ID {
		t.Errorf("wallet.ID is %v, expected %v", entityA.ID, entityB.ID)
	}
	if entityA.UserID != entityB.UserID {
		t.Errorf("wallet.UserID is %v, expected %v", entityA.UserID, entityB.UserID)
	}
	if !entityA.Balance.Equal(entityB.Balance) {
		t.Errorf("wallet.Balance is %v, expected %v", entityA.Balance, entityB.Balance)
	}
	if entityA.CreatedAt != entityB.CreatedAt {
		t.Errorf("wallet.CreatedAt is %v, expected %v", entityA.CreatedAt, entityB.CreatedAt)
	}
	if entityA.UpdatedAt != entityB.UpdatedAt {
		t.Errorf("wallet.UpdatedAt is %v, expected %v", entityA.UpdatedAt, entityB.UpdatedAt)
	}
	if entityA.DeletedAt != entityB.DeletedAt {
		t.Errorf("wallet.DeletedAt is %v, expected %v", entityA.DeletedAt, entityB.DeletedAt)
	}
}

func TestWalletUCGetWallet_WalletExists_ReturnsWallet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedEnt := GetTestWalletEntityWithDeletedAt()
	mockRepo := mocks.NewMockWalletRepoI(mockCtrl)

	retEnt := expectedEnt
	mockRepo.EXPECT().
		GetByUserID(expectedEnt.UserID).
		Return(&retEnt, nil)

	uc := usecases.NewUseCases(&infra.Repositories{WalletRepo: mockRepo}).GetWalletUC()

	entity, created, userCreated, err := uc.GetWallet(expectedEnt.UserID)
	if err != nil {
		t.Fatal(err)
	}

	CheckWalletsAreEquals(t, *entity, expectedEnt)

	if created {
		t.Errorf("the wallet was created, expected as not created")
	}
	if userCreated {
		t.Errorf("the user was created, expected as not created")
	}
}
func TestWalletUCGetWallet_WalletDoesNotExists_WalletCreated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedEnt := GetTestWalletEntity()
	expectedEnt.Balance = domain.NewMoneyFromString("0.0")
	expectedUserEnt := GetTestUserEntity()

	mockRepo := mocks.NewMockWalletRepoI(mockCtrl)
	mockUserRepo := mocks.NewMockUserRepoI(mockCtrl)

	mockRepo.EXPECT().
		GetByUserID(expectedEnt.UserID).
		Return(nil, nil)

	retUserEnt := expectedUserEnt
	mockUserRepo.EXPECT().
		GetByID(expectedUserEnt.ID).
		Return(&retUserEnt, nil)

	preInsertEnt := domain.Wallet{
		UserID:  expectedEnt.UserID,
		Balance: domain.NewMoneyFromString("0.0"),
	}

	retEnt := expectedEnt
	mockRepo.EXPECT().
		Insert(preInsertEnt).
		Return(&retEnt, nil)

	uc := usecases.NewUseCases(&infra.Repositories{UserRepo: mockUserRepo, WalletRepo: mockRepo}).GetWalletUC()

	entity, created, _, err := uc.GetWallet(expectedEnt.UserID)
	if err != nil {
		t.Fatal(err)
	}

	CheckWalletsAreEquals(t, *entity, expectedEnt)

	if !created {
		t.Errorf("the wallet not created, expected as created")
	}
}

func TestWalletUCGetWallet_BetweenGetByUserIDAndInsertCalls_ReturnsWallet(t *testing.T) {
	// This test the event of a wallet inserted between the GetByUserID and Insert calls.
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedEnt := GetTestWalletEntity()
	expectedEnt.Balance = domain.NewMoneyFromString("0.0")
	expectedUserEnt := GetTestUserEntity()

	mockRepo := mocks.NewMockWalletRepoI(mockCtrl)
	mockUserRepo := mocks.NewMockUserRepoI(mockCtrl)

	mockRepo.EXPECT().
		GetByUserID(expectedEnt.UserID).
		Return(nil, nil)

	returnedUseEnt := expectedUserEnt
	mockUserRepo.EXPECT().
		GetByID(expectedUserEnt.ID).
		Return(&returnedUseEnt, nil)

	preInsertEnt := domain.Wallet{
		UserID:  expectedEnt.UserID,
		Balance: domain.NewMoneyFromString("0.0"),
	}
	mockRepo.EXPECT().
		Insert(preInsertEnt).
		Return(nil, fmt.Errorf("%w: User ID %d", infra.ErrWalletAlreadyExists, expectedEnt.UserID))

	returnedEnt := expectedEnt
	mockRepo.EXPECT().
		GetByUserID(expectedEnt.UserID).
		Return(&returnedEnt, nil)

	uc := usecases.NewUseCases(&infra.Repositories{UserRepo: mockUserRepo, WalletRepo: mockRepo}).GetWalletUC()

	entity, created, _, err := uc.GetWallet(expectedEnt.UserID)
	if err != nil {
		t.Fatal(err)
	}

	CheckWalletsAreEquals(t, *entity, expectedEnt)

	if created {
		t.Errorf("the wallet created, expected as not created")
	}
}

func TestWalletUCIncBalance(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	credit := domain.NewMoneyFromString("999.999999")
	balance := domain.NewMoneyFromString("1000.999999")

	expectedEnt := GetTestWalletEntity()
	expectedEnt.Balance = balance.Add(credit)

	mockRepo := mocks.NewMockWalletRepoI(mockCtrl)

	returnedEnt := expectedEnt
	mockRepo.EXPECT().
		IncBalanceByUserID(expectedEnt.UserID, credit).
		Return(&returnedEnt, nil)

	uc := usecases.NewUseCases(&infra.Repositories{WalletRepo: mockRepo}).GetWalletUC()

	entity, err := uc.IncBalance(expectedEnt.UserID, credit)
	if err != nil {
		t.Fatal(err)
	}
	CheckWalletsAreEquals(t, *entity, expectedEnt)
}
