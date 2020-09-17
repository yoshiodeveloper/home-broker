package postgresql_test

import (
	"errors"
	"home-broker/domain"
	"home-broker/infra"
	"home-broker/infra/postgresql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"gorm.io/gorm"
)

var (
	walletBaseTime = time.Date(2020, time.Month(2), 0, 0, 0, 0, 0, time.UTC)
)

func GetTestWalletRepo() postgresql.WalletRepo {
	dbClient := GetTestDBClient()
	return *postgresql.NewWalletRepo(&dbClient)
}

func GetMockedWalletRepo() (postgresql.WalletRepo, sqlmock.Sqlmock) {
	dbClient, mock := GetMockedDBClient()
	return *postgresql.NewWalletRepo(&dbClient), mock
}

func GetTestWalletEntity() domain.Wallet {
	user := GetTestUserEntity()
	return domain.Wallet{
		ID:        domain.WalletID(9999),
		UserID:    user.ID,
		Balance:   domain.NewMoneyFromFloatString("999999999.999999"),
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

func GetTestWalletGORM() postgresql.WalletGORM {
	entity := GetTestWalletEntity()
	return postgresql.WalletGORM{
		ID:        int64(entity.ID),
		UserID:    int64(entity.UserID),
		Balance:   entity.Balance,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: gorm.DeletedAt{Time: time.Time{}, Valid: false},
	}
}

func GetTestWalletGORMWithDeletedAt() postgresql.WalletGORM {
	modelGORM := GetTestWalletGORM()
	entity := GetTestWalletEntityWithDeletedAt()
	modelGORM.DeletedAt = gorm.DeletedAt{Time: entity.DeletedAt, Valid: true}
	return modelGORM
}

func CheckWalletsAreEquals(t *testing.T, entityA domain.Wallet, entityB domain.Wallet) {
	if entityA.ID != entityB.ID {
		t.Errorf("wallet.ID is %v, expected %v", entityA.ID, entityB.ID)
	}
	if entityA.UserID != entityB.UserID {
		t.Errorf("wallet.UserID is %v, expected %v", entityA.UserID, entityB.UserID)
	}
	if entityA.Balance != entityB.Balance {
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

func CheckWalletGORMEntityAreEquals(t *testing.T, modelGORM postgresql.WalletGORM, entity domain.Wallet) {
	if modelGORM.ID != int64(entity.ID) {
		t.Errorf("modelGORM.ID is %v, expected %v", modelGORM.ID, int64(entity.ID))
	}
	if modelGORM.UserID != int64(entity.UserID) {
		t.Errorf("modelGORM.UserID is %v, expected %v", modelGORM.UserID, int64(entity.UserID))
	}
	if entity.Balance != modelGORM.Balance {
		t.Errorf("modelGORM.Balance is %v, expected %v", modelGORM.Balance, entity.Balance)
	}
	if modelGORM.CreatedAt != entity.CreatedAt {
		t.Errorf("modelGORM.CreatedAt is %v, expected %v", modelGORM.CreatedAt, entity.CreatedAt)
	}
	if modelGORM.UpdatedAt != entity.UpdatedAt {
		t.Errorf("modelGORM.UpdatedAt %v, expected %v", modelGORM.UpdatedAt, entity.UpdatedAt)
	}
	if modelGORM.DeletedAt.Time != entity.DeletedAt {
		t.Errorf("modelGORM.DeletedAt is %v, expected %v", modelGORM.DeletedAt.Time, entity.DeletedAt)
	}
}

func CheckWalletEntityGORMAreEquals(t *testing.T, entity domain.Wallet, modelGORM postgresql.WalletGORM) {
	if int64(entity.ID) != modelGORM.ID {
		t.Errorf("wallet.ID is %v, expected %v", int64(entity.ID), modelGORM.ID)
	}
	if int64(entity.UserID) != modelGORM.UserID {
		t.Errorf("wallet.UserID is %v, expected %v", int64(entity.UserID), modelGORM.UserID)
	}
	if entity.Balance != modelGORM.Balance {
		t.Errorf("wallet.Balance is %v, expected %v", entity.Balance, modelGORM.Balance)
	}
	if entity.CreatedAt != modelGORM.CreatedAt {
		t.Errorf("user.CreatedAt is %v, expected %v", entity.CreatedAt, modelGORM.CreatedAt)
	}
	if entity.UpdatedAt != modelGORM.UpdatedAt {
		t.Errorf("user.UpdatedAt %v, expected %v", entity.UpdatedAt, modelGORM.UpdatedAt)
	}
	if entity.DeletedAt != modelGORM.DeletedAt.Time {
		t.Errorf("user.DeletedAt is %v, expected %v", entity.DeletedAt, modelGORM.DeletedAt)
	}
}

func TestWalletRepoToWalletEntity(t *testing.T) {
	modelGORM := GetTestWalletGORMWithDeletedAt()
	repo := GetTestWalletRepo()
	entity := repo.ToEntity(&modelGORM)
	CheckWalletEntityGORMAreEquals(t, entity, modelGORM)
}

func TestWalletRepoToWalletEntity_DeletedAtFieldIsNull_WalletDeletedAtIsZero(t *testing.T) {
	modelGORM := GetTestWalletGORM()
	repo := GetTestWalletRepo()
	entity := repo.ToEntity(&modelGORM)
	CheckWalletEntityGORMAreEquals(t, entity, modelGORM)
	if !entity.DeletedAt.IsZero() {
		t.Errorf("entity.DeletedAt is %v, expected as zero", entity.DeletedAt)
	}
}

func TestWalletRepoToGORMModel(t *testing.T) {
	repo := GetTestWalletRepo()
	entity := GetTestWalletEntityWithDeletedAt()
	modelGORM := repo.ToGORMModel(&entity)
	CheckWalletGORMEntityAreEquals(t, modelGORM, entity)
}

func TestWalletRepoToGORMModel_DeletedAtFieldIsNull_WalletDeletedAtIsZero(t *testing.T) {
	repo := GetTestWalletRepo()
	entity := GetTestWalletEntity()
	entity.Balance = domain.NewMoneyFromFloatString("999999999.999999")
	modelGORM := repo.ToGORMModel(&entity)
	CheckWalletGORMEntityAreEquals(t, modelGORM, entity)
}

func TestWalletRepoGetByUserID_WalletWithUserIDDoesNotExist_NoWalletReturned(t *testing.T) {
	repo, mock := GetMockedWalletRepo()
	userID := int64(9999)
	mock.ExpectQuery(`SELECT \* FROM "wallet" WHERE "wallet"\."user_id" = \$1 AND "wallet"\."deleted_at" IS NULL LIMIT 1`).
		WithArgs(userID).
		WillReturnRows(mock.NewRowsWithColumnDefinition())
	entity, err := repo.GetByUserID(domain.UserID(userID))
	if err != nil {
		t.Fatal(err)
	}
	if entity != nil {
		t.Error("wallet found, expected as not found")
	}
}

func TestWalletRepoGetByUserID_WalletWithUserIDExists_WalletReturned(t *testing.T) {
	repo, mock := GetMockedWalletRepo()
	expectedEntity := GetTestWalletEntity()
	columns := []string{"id", "user_id", "balance", "created_at", "updated_at", "deleted_at"}
	mock.ExpectQuery(`SELECT \* FROM "wallet" WHERE "wallet"\."user_id" = \$1 AND "wallet"\."deleted_at" IS NULL LIMIT 1`).
		WithArgs(int64(expectedEntity.UserID)).
		WillReturnRows(mock.NewRows(columns).
			AddRow(
				int64(expectedEntity.ID), int64(expectedEntity.UserID), expectedEntity.Balance,
				expectedEntity.CreatedAt, expectedEntity.UpdatedAt, nil))
	entity, err := repo.GetByUserID(expectedEntity.UserID)
	if err != nil {
		t.Fatal(err)
	}
	if entity == nil {
		t.Error("wallet not found, expected as found")
	}
	CheckWalletsAreEquals(t, *entity, expectedEntity)
}

func TestWalletRepoInsert_UserExists_NewWalletInserted(t *testing.T) {
	repo, mock := GetMockedWalletRepo()
	expectedEntity := GetTestWalletEntity()
	expectedEntity.ID = domain.WalletID(9999) // future ID
	mock.ExpectQuery(`INSERT INTO "wallet" \("created_at","updated_at","deleted_at","user_id","balance"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			int64(expectedEntity.UserID),
			expectedEntity.Balance,
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(expectedEntity.ID)))

	entity := GetTestWalletEntity()
	entity.ID = 0
	newEntity, err := repo.Insert(entity)
	if err != nil {
		t.Fatal(err)
	}
	CheckWalletsAreEquals(t, *newEntity, expectedEntity)
}

func TestWalletRepoInsert_WalletIDExists_WalletIDAlreadyExistsError(t *testing.T) {
	repo, mock := GetMockedWalletRepo()
	expectedEntity := GetTestWalletEntity()
	expectedEntity.ID = domain.WalletID(9999) // future ID

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "wallet" \("created_at","updated_at","deleted_at","user_id","balance","id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			int64(expectedEntity.UserID),
			expectedEntity.Balance,
			int64(expectedEntity.ID),
		).WillReturnError(errors.New(`ERROR: duplicate key value violates unique constraint "wallet_pkey" (SQLSTATE 23505)`))
	mock.ExpectRollback()

	entity := GetTestWalletEntity()
	entity.ID = expectedEntity.ID
	newEntity, err := repo.Insert(entity)

	if err == nil {
		t.Error("wallet created, expected ErrWalletAlreadyExists")
	}
	if !errors.Is(err, infra.ErrWalletAlreadyExists) {
		t.Errorf("expected ErrWalletAlreadyExists, received \"%v\"", err)
	}
	if newEntity != nil {
		t.Errorf("wallet received \"%v\", expected nil", newEntity)
	}
}

func TestWalletRepoInsert_WalletUserIDExists_WalletUserIDAlreadyExistsError(t *testing.T) {
	repo, mock := GetMockedWalletRepo()
	expectedEntity := GetTestWalletEntity()
	expectedEntity.ID = domain.WalletID(9999) // future ID

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "wallet" \("created_at","updated_at","deleted_at","user_id","balance","id"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			int64(expectedEntity.UserID),
			expectedEntity.Balance,
		).WillReturnError(errors.New(`ERROR: duplicate key value violates unique constraint "wallet_user_id_key" (SQLSTATE 23505)`))
	mock.ExpectRollback()

	entity := GetTestWalletEntity()
	entity.ID = 0
	newEntity, err := repo.Insert(entity)

	if err == nil {
		t.Error("wallet created, expected ErrWalletAlreadyExists")
	}
	if !errors.Is(err, infra.ErrWalletAlreadyExists) {
		t.Errorf("expected ErrWalletAlreadyExists, received \"%v\"", err)
	}
	if newEntity != nil {
		t.Errorf("wallet received \"%v\", expected nil", newEntity)
	}
}

func TestWalletRepoInsert_WalletUserIDDoesNotExist_UserDoesNotExistError(t *testing.T) {
	repo, mock := GetMockedWalletRepo()
	expectedEntity := GetTestWalletEntity()
	expectedEntity.ID = domain.WalletID(9999) // future ID

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "wallet" \("created_at","updated_at","deleted_at","user_id","balance","id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			int64(expectedEntity.UserID),
			expectedEntity.Balance,
			int64(expectedEntity.ID),
		).WillReturnError(errors.New(`ERROR: insert or update on table "wallet" violates foreign key constraint "fk_wallet_user" (SQLSTATE 23503)`))
	mock.ExpectRollback()

	entity := GetTestWalletEntity()
	entity.ID = expectedEntity.ID
	newEntity, err := repo.Insert(entity)

	if err == nil {
		t.Error("wallet created, expected ErrUserDoesNotExist")
	}
	if !errors.Is(err, infra.ErrUserDoesNotExist) {
		t.Errorf("expected ErrUserDoesNotExist, received \"%v\"", err)
	}
	if newEntity != nil {
		t.Errorf("wallet received \"%v\", expected nil", newEntity)
	}
}

func TestWalletRepoIncBalanceByUserID(t *testing.T) {
	repo, mock := GetMockedWalletRepo()
	expectedEnt := GetTestWalletEntity()

	value := domain.NewMoneyFromFloatString("999.999999")
	balance := domain.NewMoneyFromFloatString("1000.999999")
	expectedEnt.Balance += balance

	mock.ExpectExec(`UPDATE "wallet" SET .+ WHERE.+"user_id"\s*=\s*\$3`).
		WithArgs(
			value,
			sqlmock.AnyArg(),
			int64(expectedEnt.UserID),
		).WillReturnResult(sqlmock.NewResult(0, 1))

	columns := []string{"id", "user_id", "balance", "created_at", "updated_at", "deleted_at"}

	mock.ExpectQuery(`SELECT \* FROM "wallet".+WHERE.+"user_id"\s*=\s*\$1.*`).
		WithArgs(int64(expectedEnt.UserID)).
		WillReturnRows(mock.NewRows(columns).
			AddRow(
				int64(expectedEnt.ID), int64(expectedEnt.UserID), expectedEnt.Balance,
				expectedEnt.CreatedAt, expectedEnt.UpdatedAt, nil))

	entity, err := repo.IncBalanceByUserID(expectedEnt.UserID, value)
	if err != nil {
		t.Fatal(err)
	}
	CheckWalletsAreEquals(t, *entity, expectedEnt)
}

func TestWalletRepoIncBalanceByUserID_WalletDoesNotExist_NoWalletReturned(t *testing.T) {
	repo, mock := GetMockedWalletRepo()
	expectedEnt := GetTestWalletEntity()

	value := domain.NewMoneyFromFloatString("999.999999")
	balance := domain.NewMoneyFromFloatString("1000.999999")
	expectedEnt.Balance += balance

	mock.ExpectExec(`UPDATE "wallet" SET .+ WHERE.+"user_id"\s*=\s*\$3`).
		WithArgs(
			value,
			sqlmock.AnyArg(),
			int64(expectedEnt.UserID),
		).WillReturnResult(sqlmock.NewResult(0, 0))

	columns := []string{"id", "user_id", "balance", "created_at", "updated_at", "deleted_at"}

	mock.ExpectQuery(`SELECT \* FROM "wallet".+WHERE.+"user_id"\s*=\s*\$1.*`).
		WithArgs(int64(expectedEnt.UserID)).
		WillReturnRows(mock.NewRows(columns))

	entity, err := repo.IncBalanceByUserID(expectedEnt.UserID, value)
	if err != nil {
		t.Fatal(err)
	}
	if entity != nil {
		t.Errorf("wallet exists, expected as does not exist")
	}
}
