package postgresql_test

import (
	"errors"
	"home-broker/money"
	postgresqltests "home-broker/tests/postgresql"
	walletstests "home-broker/tests/wallets"
	"home-broker/users"
	"home-broker/wallets"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestToWalletEntity(t *testing.T) {
	model := postgresqltests.GetWalletModelWithDeletedAt()
	db := postgresqltests.GetWalletDB()
	entity := db.ToEntity(model)
	err := postgresqltests.CheckWalletsEntityModel(entity, model)
	if err != nil {
		t.Error(err)
	}
}

func TestToWalletEntity_DeletedAtFieldIsNull_WalletDeletedAtIsZero(t *testing.T) {
	model := postgresqltests.GetWalletModel()
	db := postgresqltests.GetWalletDB()
	entity := db.ToEntity(model)

	err := postgresqltests.CheckWalletsEntityModel(entity, model)
	if err != nil {
		t.Error(err)
	}
	if !entity.DeletedAt.IsZero() {
		t.Errorf("entity.DeletedAt is %v, expected as zero", entity.DeletedAt)
	}
}

func TestToModel_DeletedAtFieldIsNull_WalletDeletedAtIsZero(t *testing.T) {
	db := postgresqltests.GetWalletDB()
	entity := walletstests.GetWallet()
	model := db.ToModel(entity)
	err := postgresqltests.CheckWalletsModelEntity(model, entity)
	if err != nil {
		t.Error(err)
	}
}

func TestToModel(t *testing.T) {
	db := postgresqltests.GetWalletDB()
	entity := walletstests.GetWalletWithDeletedAt()
	model := db.ToModel(entity)
	err := postgresqltests.CheckWalletsModelEntity(model, entity)
	if err != nil {
		t.Error(err)
	}
}

func TestGetByUserID_WalletWithUserIDDoesNotExist_NoWalletReturned(t *testing.T) {
	db, mock, err := postgresqltests.GetMockedWalletDB()
	if err != nil {
		t.Error(err)
	}
	userID := users.UserID(9999)
	mock.ExpectQuery(`SELECT \* FROM "wallet" WHERE "wallet"\."user_id" = \$1 AND "wallet"\."deleted_at" IS NULL LIMIT 1`).
		WithArgs(userID).
		WillReturnRows(mock.NewRowsWithColumnDefinition())
	entity, err := db.GetByUserID(userID)
	if err != nil {
		t.Fatal(err)
	}
	if entity != nil {
		t.Error("wallet found, expected as not found")
	}
}

func TestGetByUserID_WalletWithUserIDExists_WalletReturned(t *testing.T) {
	db, mock, err := postgresqltests.GetMockedWalletDB()
	if err != nil {
		t.Error(err)
	}
	expectedEntity := walletstests.GetWallet()

	columns := []string{"id", "user_id", "balance", "created_at", "updated_at", "deleted_at"}
	mock.ExpectQuery(`SELECT \* FROM "wallet" WHERE "wallet"\."user_id" = \$1 AND "wallet"\."deleted_at" IS NULL LIMIT 1`).
		WithArgs(int64(expectedEntity.UserID)).
		WillReturnRows(mock.NewRows(columns).
			AddRow(
				expectedEntity.ID, expectedEntity.UserID, expectedEntity.Balance,
				expectedEntity.CreatedAt, expectedEntity.UpdatedAt, nil))

	entity, err := db.GetByUserID(expectedEntity.UserID)
	if err != nil {
		t.Error(err)
	}
	if entity == nil {
		t.Error("wallet not found, expected as found")
	}

	err = walletstests.CheckWallets(*entity, expectedEntity)
	if err != nil {
		t.Error(err)
	}
}

func TestInsert_UserExists_NewWalletInserted(t *testing.T) {
	db, mock, err := postgresqltests.GetMockedWalletDB()
	if err != nil {
		t.Error(err)
	}
	expectedEntity := walletstests.GetWallet()
	expectedEntity.ID = wallets.WalletID(9999) // future ID
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "wallet" \("created_at","updated_at","deleted_at","user_id","balance"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			int64(expectedEntity.UserID),
			expectedEntity.Balance,
		).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedEntity.ID))
	mock.ExpectCommit()

	entity := walletstests.GetWallet()

	entity.ID = 0
	newEntity, err := db.Insert(entity)
	if err != nil {
		t.Error(err)
	}
	if newEntity == nil {
		t.Error("newEntity is nil")
	}

	err = walletstests.CheckWallets(*newEntity, expectedEntity)
	if err != nil {
		t.Error(err)
	}
}
func TestInsert_WalletIDExists_WalletIDAlreadyExistsError(t *testing.T) {
	db, mock, err := postgresqltests.GetMockedWalletDB()
	if err != nil {
		t.Error(err)
	}
	expectedEntity := walletstests.GetWallet()
	expectedEntity.ID = wallets.WalletID(9999) // future ID
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "wallet" \("created_at","updated_at","deleted_at","user_id","balance","id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			expectedEntity.UserID,
			expectedEntity.Balance,
			expectedEntity.ID,
		).WillReturnError(errors.New(`ERROR: duplicate key value violates unique constraint "wallet_pkey" (SQLSTATE 23505)`))
	mock.ExpectRollback()

	entity := walletstests.GetWallet()

	entity.ID = expectedEntity.ID
	newEntity, err := db.Insert(entity)

	if err == nil {
		t.Error("wallet created, expected ErrWalletAlreadyExists")
	}
	if !errors.Is(err, wallets.ErrWalletAlreadyExists) {
		t.Errorf("expected ErrWalletAlreadyExists, received \"%v\"", err)
	}
	if newEntity != nil {
		t.Errorf("wallet received \"%v\", expected nil", newEntity)
	}
}

func TestInsert_WalletUserIDExists_WalletUserIDAlreadyExistsError(t *testing.T) {
	db, mock, err := postgresqltests.GetMockedWalletDB()
	if err != nil {
		t.Error(err)
	}
	expectedEntity := walletstests.GetWallet()
	expectedEntity.ID = wallets.WalletID(9999) // future ID
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "wallet" \("created_at","updated_at","deleted_at","user_id","balance","id"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			expectedEntity.UserID,
			expectedEntity.Balance,
		).WillReturnError(errors.New(`ERROR: duplicate key value violates unique constraint "wallet_user_id_key" (SQLSTATE 23505)`))
	mock.ExpectRollback()

	entity := walletstests.GetWallet()

	entity.ID = 0
	newEntity, err := db.Insert(entity)

	if err == nil {
		t.Error("wallet created, expected ErrWalletAlreadyExists")
	}
	if !errors.Is(err, wallets.ErrWalletAlreadyExists) {
		t.Errorf("expected ErrWalletAlreadyExists, received \"%v\"", err)
	}
	if newEntity != nil {
		t.Errorf("wallet received \"%v\", expected nil", newEntity)
	}
}

func TestInsert_WalletUserIDDoesNotExist_UserDoesNotExistError(t *testing.T) {
	db, mock, err := postgresqltests.GetMockedWalletDB()
	if err != nil {
		t.Error(err)
	}
	expectedEntity := walletstests.GetWallet()
	expectedEntity.ID = wallets.WalletID(9999) // future ID

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "wallet" \("created_at","updated_at","deleted_at","user_id","balance","id"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			expectedEntity.UserID,
			expectedEntity.Balance,
			expectedEntity.ID,
		).WillReturnError(errors.New(`ERROR: insert or update on table "wallet" violates foreign key constraint "fk_wallet_user" (SQLSTATE 23503)`))
	mock.ExpectRollback()

	entity := walletstests.GetWallet()

	entity.ID = expectedEntity.ID
	newEntity, err := db.Insert(entity)

	if err == nil {
		t.Error("wallet created, expected ErrUserDoesNotExist")
	}
	if !errors.Is(err, users.ErrUserDoesNotExist) {
		t.Errorf("expected ErrUserDoesNotExist, received \"%v\"", err)
	}
	if newEntity != nil {
		t.Errorf("wallet received \"%v\", expected nil", newEntity)
	}
}

func TestIncBalanceByUserID(t *testing.T) {
	db, mock, err := postgresqltests.GetMockedWalletDB()
	if err != nil {
		t.Error(err)
	}
	expectedEntity := walletstests.GetWallet()

	value, err := money.NewMoneyFromFloatString("999.999999")
	if err != nil {
		t.Error(err)
	}
	balance, err := money.NewMoneyFromFloatString("1000.999999")
	if err != nil {
		t.Error(err)
	}

	expectedEntity.Balance += balance

	t.Run("WalletExists_WalletRetuned", func(t *testing.T) {
		mock.ExpectExec(`UPDATE "wallet" SET .+ WHERE.+"user_id"\s*=\s*\$3`).
			WithArgs(
				value,
				sqlmock.AnyArg(),
				expectedEntity.UserID,
			).WillReturnResult(sqlmock.NewResult(0, 1))

		columns := []string{"id", "user_id", "balance", "created_at", "updated_at", "deleted_at"}

		mock.ExpectQuery(`SELECT \* FROM "wallet".+WHERE.+"user_id"\s*=\s*\$1.*`).
			WithArgs(expectedEntity.UserID).
			WillReturnRows(mock.NewRows(columns).
				AddRow(
					expectedEntity.ID, expectedEntity.UserID, expectedEntity.Balance,
					expectedEntity.CreatedAt, expectedEntity.UpdatedAt, nil))

		entity, err := db.IncBalanceByUserID(expectedEntity.UserID, value)
		if err != nil {
			t.Error(err)
		}

		err = walletstests.CheckWallets(*entity, expectedEntity)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("WalletExists_NoWalletReturned", func(t *testing.T) {
		mock.ExpectExec(`UPDATE "wallet" SET .+ WHERE.+"user_id"\s*=\s*\$3`).
			WithArgs(
				value,
				sqlmock.AnyArg(),
				expectedEntity.UserID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		columns := []string{"id", "user_id", "balance", "created_at", "updated_at", "deleted_at"}

		mock.ExpectQuery(`SELECT \* FROM "wallet".+WHERE.+"user_id"\s*=\s*\$1.*`).
			WithArgs(expectedEntity.UserID).
			WillReturnRows(mock.NewRows(columns))

		entity, err := db.IncBalanceByUserID(expectedEntity.UserID, value)
		if err != nil {
			t.Error(err)
		}
		if entity != nil {
			t.Errorf("wallet exists, expected as does not exist")
		}
	})
}
