package postgresql_test

import (
	"errors"
	postgresqltests "home-broker/tests/postgresql"
	userstests "home-broker/tests/users"
	"home-broker/users"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestUserDBToUserEntity(t *testing.T) {
	model := postgresqltests.GetUserModel()
	db := postgresqltests.GetUserDB()
	entity := db.ToEntity(model)
	if entity.ID != model.ID {
		t.Errorf("user.ID is %v, expected %v", int64(entity.ID), model.ID)
	}
	if !entity.CreatedAt.Equal(model.CreatedAt) {
		t.Errorf("user.CreatedAt is %v, expected %v", entity.CreatedAt, model.CreatedAt)
	}
	if !entity.UpdatedAt.Equal(model.UpdatedAt) {
		t.Errorf("user.UpdatedAt %v, expected %v", entity.UpdatedAt, model.UpdatedAt)
	}
	if !entity.DeletedAt.Equal(model.DeletedAt.Time) {
		t.Errorf("user.DeletedAt is %v, expected %v", entity.DeletedAt, model.DeletedAt.Time)
	}
}

func TestUserDBToUserEntity_DeletedAtFieldIsNull_UserDeletedAtIsZero(t *testing.T) {
	model := postgresqltests.GetUserModel()
	model.DeletedAt = gorm.DeletedAt{Time: userstests.BaseTime.Add(time.Hour * 3), Valid: false}
	db := postgresqltests.GetUserDB()
	user := db.ToEntity(model)
	if !user.DeletedAt.IsZero() {
		t.Errorf("user.DeletedAt is %v, expected as zero", user.DeletedAt)
	}
}

func TestUserDBGetByID_UserIDDoesNotExist_NoUserReturned(t *testing.T) {
	db, mock, err := postgresqltests.GetMockedUserDB()
	if err != nil {
		t.Error(err)
	}
	userID := users.UserID(99)
	mock.ExpectQuery(`SELECT \* FROM "user" WHERE "user"\."id" = \$1 AND "user"\."deleted_at" IS NULL LIMIT 1`).
		WithArgs(userID).
		WillReturnRows(mock.NewRowsWithColumnDefinition())

	entity, err := db.GetByID(userID)
	if err != nil {
		t.Fatal(err)
	}
	if entity != nil {
		t.Errorf("user found (%v), expected as not found", entity)
	}
}

func TestUserDBGetByID_UserIDExists_UserMustBeReturned(t *testing.T) {
	db, mock, err := postgresqltests.GetMockedUserDB()
	if err != nil {
		t.Error(err)
	}
	expectedEntity := userstests.GetEntity()
	columns := []string{"id", "created_at", "updated_at", "deleted_at"}
	mock.ExpectQuery(`SELECT \* FROM "user" WHERE "user"\."id" = \$1 AND "user"\."deleted_at" IS NULL LIMIT 1`).
		WithArgs(int64(expectedEntity.ID)).
		WillReturnRows(mock.NewRows(columns).AddRow(int64(expectedEntity.ID), expectedEntity.CreatedAt, expectedEntity.UpdatedAt, nil))
	entity, err := db.GetByID(expectedEntity.ID)
	if err != nil {
		t.Fatal(err)
	}
	if entity == nil {
		t.Error("user not found, expected as found")
	}
	if !reflect.DeepEqual(*entity, expectedEntity) {
		t.Errorf("received %v, expected %v", *entity, expectedEntity)
	}
}

func TestUserDBInsert_UserIDDoesNotExist_NewUserMustBeReturned(t *testing.T) {
	db, mock, err := postgresqltests.GetMockedUserDB()
	if err != nil {
		t.Error(err)
	}
	userID := users.UserID(999)
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "user" \("id","created_at","updated_at","deleted_at"\) VALUES \(\$1,\$2,\$3,\$4\)`).
		WithArgs(userID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnResult(sqlmock.NewResult(int64(userID), 1))
	mock.ExpectCommit()
	entity, err := db.Insert(users.User{ID: userID})
	if err != nil {
		if errors.Is(err, users.ErrUserAlreadyExists) {
			t.Error("user already exists, expected as does not exist")
		} else {
			t.Fatal(err)
		}
	}
	if entity.ID != users.UserID(userID) {
		t.Errorf("user.ID is %v, expected %v", entity.ID, userID)
	}
	if entity.CreatedAt.IsZero() {
		t.Errorf("user.CreatedAt is zero, expected time.Time")
	}
	if entity.UpdatedAt.IsZero() {
		t.Errorf("user.UpdatedAt is zero, expected time.Time")
	}
	if !entity.DeletedAt.IsZero() {
		t.Errorf("user.DeletedAt is not zero, expected zero")
	}
}

func TestUserDBInsert_UserIDExists_AlreadExistsReturned(t *testing.T) {
	db, mock, err := postgresqltests.GetMockedUserDB()
	if err != nil {
		t.Error(err)
	}
	userID := users.UserID(999)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "user" \("id","created_at","updated_at","deleted_at"\) VALUES \(\$1,\$2,\$3,\$4\)`).
		WithArgs(userID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))
	mock.ExpectRollback()

	entity, err := db.Insert(users.User{ID: userID})
	if !errors.Is(err, users.ErrUserAlreadyExists) {
		if err == nil {
			t.Error("user does not exists, expected already exists error")
		} else {
			t.Fatal(err)
		}
	}
	if entity != nil {
		t.Errorf("user returned (%v), expected nil", entity)
	}
}
