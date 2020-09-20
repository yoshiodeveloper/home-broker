package users_test

import (
	"fmt"
	testusers "home-broker/tests/users"
	"home-broker/tests/users/mocks"
	"home-broker/users"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUserUseCases(t *testing.T) {
	expectedEnt := testusers.GetEntity()
	retEnt := expectedEnt

	t.Run("UserExists_ReturnsUser", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockDB := mocks.NewMockUserDBInterface(mockCtrl)
		uc := users.NewUserUseCases(mockDB)

		mockDB.EXPECT().
			GetByID(expectedEnt.ID).
			Return(&retEnt, nil)

		entity, created, err := uc.GetUser(expectedEnt.ID)
		if err != nil {
			t.Fatal(err)
		}
		err = testusers.CheckUsers(*entity, expectedEnt)
		if err != nil {
			t.Error(err)
		}
		if created {
			t.Errorf("the user was created, expected as not created")
		}
	})

	t.Run("UserDoesNotExist_ReturnsNewUser", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockDB := mocks.NewMockUserDBInterface(mockCtrl)
		uc := users.NewUserUseCases(mockDB)

		mockDB.EXPECT().
			GetByID(expectedEnt.ID).
			Return(nil, nil)

		mockDB.EXPECT().
			Insert(users.User{ID: expectedEnt.ID}).
			Return(&retEnt, nil)

		entity, created, err := uc.GetUser(expectedEnt.ID)
		if err != nil {
			t.Fatal(err)
		}

		err = testusers.CheckUsers(*entity, expectedEnt)
		if err != nil {
			t.Error(err)
		}
		if created {
			t.Errorf("the user was not created, expected as created")
		}
	})

	// This test the event of an user inserted between the GetByID and Insert calls.
	t.Run("BetweenGetByIDAndInsertCalls_ReturnsUser", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockDB := mocks.NewMockUserDBInterface(mockCtrl)
		uc := users.NewUserUseCases(mockDB)

		mockDB.EXPECT().
			GetByID(expectedEnt.ID).
			Return(nil, nil)

		mockDB.EXPECT().
			Insert(users.User{ID: expectedEnt.ID}).
			Return(nil, fmt.Errorf("%w: ID %d", users.ErrUserAlreadyExists, expectedEnt.ID))

		mockDB.EXPECT().
			GetByID(expectedEnt.ID).
			Return(&retEnt, nil)

		entity, created, err := uc.GetUser(expectedEnt.ID)
		if err != nil {
			t.Fatal(err)
		}
		if entity == nil {
			t.Errorf("entity is nil")
		}

		err = testusers.CheckUsers(*entity, expectedEnt)
		if err != nil {
			t.Error(err)
		}
		if created {
			t.Errorf("the user was created, expected as not created")
		}
	})
}
