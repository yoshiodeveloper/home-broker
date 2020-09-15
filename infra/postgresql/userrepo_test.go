package postgresql_test

import (
	"errors"
	"home-broker/domain"
	"home-broker/infra"
	"home-broker/infra/postgresql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"gorm.io/gorm"
)

var (
	userBaseTime = time.Date(2020, time.Month(1), 0, 0, 0, 0, 0, time.UTC)
)

func GetTestUserRepo() postgresql.UserRepo {
	dbClient := GetTestDBClient()
	return *postgresql.NewUserRepo(&dbClient)
}

func GetMockedUserRepo() (postgresql.UserRepo, sqlmock.Sqlmock) {
	dbClient, mock := GetMockedDBClient()
	return *postgresql.NewUserRepo(&dbClient), mock
}

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

func GetTestUserGORM() postgresql.UserGORM {
	entity := GetTestUserEntity()
	return postgresql.UserGORM{
		ID:        int64(entity.ID),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: gorm.DeletedAt{Time: time.Time{}, Valid: false},
	}
}

func GetTestUserGORMWithDeletedAt() postgresql.UserGORM {
	modelGORM := GetTestUserGORM()
	entity := GetTestUserEntityWithDeletedAt()
	modelGORM.DeletedAt = gorm.DeletedAt{Time: entity.DeletedAt, Valid: true}
	return modelGORM
}

func TestUserRepoToUserEntity(t *testing.T) {
	userGORM := GetTestUserGORMWithDeletedAt()
	userRepo := GetTestUserRepo()
	user := userRepo.ToEntity(&userGORM)
	if int64(user.ID) != userGORM.ID {
		t.Errorf("user.ID is %v, expected %v", int64(user.ID), userGORM.ID)
	}
	if user.CreatedAt != userGORM.CreatedAt {
		t.Errorf("user.CreatedAt is %v, expected %v", user.CreatedAt, userGORM.CreatedAt)
	}
	if user.UpdatedAt != userGORM.UpdatedAt {
		t.Errorf("user.UpdatedAt %v, expected %v", user.UpdatedAt, userGORM.UpdatedAt)
	}
	if user.DeletedAt != userGORM.DeletedAt.Time {
		t.Errorf("user.DeletedAt is %v, expected %v", user.DeletedAt, userGORM.DeletedAt)
	}
}

func TestUserRepoToUserEntity_DeletedAtFieldIsNull_UserDeletedAtIsZero(t *testing.T) {
	userGORM := GetTestUserGORM()
	userGORM.DeletedAt = gorm.DeletedAt{Time: userBaseTime.Add(time.Hour * 3), Valid: false}
	userRepo := GetTestUserRepo()
	user := userRepo.ToEntity(&userGORM)
	if !user.DeletedAt.IsZero() {
		t.Errorf("user.DeletedAt is %v, expected as zero", user.DeletedAt)
	}
}

func TestUserRepoGetByID_UserIDDoesNotExist_NoUserReturned(t *testing.T) {
	userRepo, mock := GetMockedUserRepo()
	userID := int64(99)
	mock.ExpectQuery(`SELECT \* FROM "user" WHERE "user"\."id" = \$1 AND "user"\."deleted_at" IS NULL LIMIT 1`).
		WithArgs(userID).
		WillReturnRows(mock.NewRowsWithColumnDefinition())

	entity, err := userRepo.GetByID(domain.UserID(userID))
	if err != nil {
		t.Fatal(err)
	}
	if entity != nil {
		t.Errorf("user found (%v), expected as not found", entity)
	}
}

func TestUserRepoGetByID_UserIDExists_UserMustBeReturned(t *testing.T) {
	repo, mock := GetMockedUserRepo()
	expectedEntity := GetTestUserEntity()
	columns := []string{"id", "created_at", "updated_at", "deleted_at"}
	mock.ExpectQuery(`SELECT \* FROM "user" WHERE "user"\."id" = \$1 AND "user"\."deleted_at" IS NULL LIMIT 1`).
		WithArgs(int64(expectedEntity.ID)).
		WillReturnRows(mock.NewRows(columns).AddRow(int64(expectedEntity.ID), expectedEntity.CreatedAt, expectedEntity.UpdatedAt, nil))
	entity, err := repo.GetByID(expectedEntity.ID)
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

func TestUserRepoInsert_UserIDDoesNotExist_NewUserMustBeReturned(t *testing.T) {
	repo, mock := GetMockedUserRepo()
	userID := int64(999)
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "user" \("id","created_at","updated_at","deleted_at"\) VALUES \(\$1,\$2,\$3,\$4\)`).
		WithArgs(userID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnResult(sqlmock.NewResult(userID, 1))
	mock.ExpectCommit()
	entity, err := repo.Insert(domain.User{ID: domain.UserID(userID)})
	if err != nil {
		if errors.Is(err, infra.ErrUserAlreadyExists) {
			t.Error("user already exists, expected as does not exist")
		} else {
			t.Fatal(err)
		}
	}
	if entity.ID != domain.UserID(userID) {
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

func TestUserRepoInsert_UserIDExists_AlreadExistsReturned(t *testing.T) {
	userRepo, mock := GetMockedUserRepo()
	userID := int64(999)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "user" \("id","created_at","updated_at","deleted_at"\) VALUES \(\$1,\$2,\$3,\$4\)`).
		WithArgs(userID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnError(errors.New("duplicate key value violates unique constraint"))
	mock.ExpectRollback()

	entity, err := userRepo.Insert(domain.User{ID: domain.UserID(userID)})
	if !errors.Is(err, infra.ErrUserAlreadyExists) {
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
