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
	userBaseTime = time.Date(2020, time.Month(1), 0, 0, 0, 0, 0, time.UTC)
)

func GetTestUserEntity() domain.User {
	return domain.User{
		ID:        domain.UserID(999),
		CreatedAt: userBaseTime,
		UpdatedAt: userBaseTime.Add(time.Hour * 2),
		DeletedAt: time.Time{},
	}
}

func GetTestUserEntityWithDeletedAt() domain.User {
	entity := GetTestUserEntity()
	entity.DeletedAt = userBaseTime.Add(time.Hour * 3)
	return entity
}

func CheckUsersAreEquals(t *testing.T, entityA domain.User, entityB domain.User) {
	if entityA.ID != entityB.ID {
		t.Errorf("user.ID is %v, expected %v", entityA.ID, entityB.ID)
	}
	if entityA.CreatedAt != entityB.CreatedAt {
		t.Errorf("user.CreatedAt is %v, expected %v", entityA.CreatedAt, entityB.CreatedAt)
	}
	if entityA.UpdatedAt != entityB.UpdatedAt {
		t.Errorf("user.UpdatedAt is %v, expected %v", entityA.UpdatedAt, entityB.UpdatedAt)
	}
	if entityA.DeletedAt != entityB.DeletedAt {
		t.Errorf("user.DeletedAt is %v, expected %v", entityA.DeletedAt, entityB.DeletedAt)
	}
}

func TestUserUCGetUser_UserExists_ReturnsUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedEnt := GetTestUserEntity()
	mockRepo := mocks.NewMockUserRepoI(mockCtrl)

	mockRepo.EXPECT().
		GetByID(expectedEnt.ID).
		Return(&expectedEnt, nil)

	uc := usecases.NewUseCases(&infra.Repositories{UserRepo: mockRepo}).GetUserUC()

	entity, created, err := uc.GetUser(expectedEnt.ID)
	if err != nil {
		t.Fatal(err)
	}

	CheckUsersAreEquals(t, *entity, expectedEnt)

	if created {
		t.Errorf("the user was created, expected as not created")
	}
}

func TestUserUCGetUser_UserDoesNotExist_ReturnsCreatedUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedEnt := GetTestUserEntity()
	mockRepo := mocks.NewMockUserRepoI(mockCtrl)

	mockRepo.EXPECT().
		GetByID(expectedEnt.ID).
		Return(nil, nil)

	mockRepo.EXPECT().
		Insert(domain.User{ID: expectedEnt.ID}).
		Return(&expectedEnt, nil)

	uc := usecases.NewUseCases(&infra.Repositories{UserRepo: mockRepo}).GetUserUC()

	entity, created, err := uc.GetUser(expectedEnt.ID)
	if err != nil {
		t.Fatal(err)
	}

	CheckUsersAreEquals(t, *entity, expectedEnt)

	if created {
		t.Errorf("the user was not created, expected as created")
	}
}

func TestUserUCGetUser_BetweenGetByIDAndInsertCalls_ReturnsUser(t *testing.T) {
	// This test the event of an user inserted between the GetByID and Insert calls.
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedEnt := GetTestUserEntity()
	mockRepo := mocks.NewMockUserRepoI(mockCtrl)

	mockRepo.EXPECT().
		GetByID(expectedEnt.ID).
		Return(nil, nil)

	mockRepo.EXPECT().
		Insert(domain.User{ID: expectedEnt.ID}).
		Return(nil, fmt.Errorf("%w: ID %d", infra.ErrUserAlreadyExists, expectedEnt.ID))

	mockRepo.EXPECT().
		GetByID(expectedEnt.ID).
		Return(&expectedEnt, nil)

	uc := usecases.NewUseCases(&infra.Repositories{UserRepo: mockRepo}).GetUserUC()

	entity, created, err := uc.GetUser(expectedEnt.ID)
	if err != nil {
		t.Fatal(err)
	}

	CheckUsersAreEquals(t, *entity, expectedEnt)

	if created {
		t.Errorf("the user was created, expected as not created")
	}
}
