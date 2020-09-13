package usecases

import (
	"home-broker/domain"
	"home-broker/infra"
	"home-broker/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetUser_UserDoesNotExist_CreatesNewUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	userID := domain.UserID(99)
	newUser := domain.User{
		ID: userID,
	}
	mockRepo := mocks.NewMockUserRepo(mockCtrl)
	mockRepo.EXPECT().
		GetByID(userID).
		Return(domain.User{}, false)

	mockRepo.EXPECT().
		Insert(newUser).
		Return(domain.User{ID: userID})

	uc := NewUseCases(&infra.Repositories{UserRepo: mockRepo}).GetUserUC()

	user, created := uc.GetUser(userID)
	if user.ID != userID {
		t.Errorf("the user was not created. User ID is %v, expected %v", user.ID, userID)
	}
	if !created {
		t.Error("the user was not created. \"created\" flag is false, expected true")
	}
}
