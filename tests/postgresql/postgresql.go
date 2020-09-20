package postgresql

import (
	"fmt"
	"home-broker/core/implem/postgresql"
	testusers "home-broker/tests/users"
	tests "home-broker/tests/wallets"
	userspostgresql "home-broker/users/implem/postgresql"
	"home-broker/wallets"
	walletspostgresql "home-broker/wallets/implem/postgresql"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetDB retuns a PostgreSQL DB handler.
func GetDB() postgresql.DB {
	db := postgresql.NewDB(
		"Host",
		1234,
		"User",
		"Password",
		"DBName",
	)
	return db
}

// GetMockedDB returns a mocked PostgreSQL DB.
func GetMockedDB() (postgresql.DB, sqlmock.Sqlmock, error) {
	db := GetDB()
	mockedDB, mock, err := sqlmock.New()
	if err != nil {
		return postgresql.DB{}, nil, err
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: mockedDB}), &gorm.Config{})
	db.SetDB(gormDB)
	return db, mock, nil
}

// GetUserDB returns a PostgreSQL user DB.
func GetUserDB() userspostgresql.UserDB {
	return userspostgresql.NewUserDB(GetDB())
}

// GetMockedUserDB returns a mocked PostgreSQL user DB.
func GetMockedUserDB() (userspostgresql.UserDB, sqlmock.Sqlmock, error) {
	db, mock, err := GetMockedDB()
	if err != nil {
		return userspostgresql.UserDB{}, nil, err
	}
	return userspostgresql.NewUserDB(db), mock, nil
}

// GetUserModel returns a PostgreSQL user model.
func GetUserModel() userspostgresql.UserModel {
	entity := testusers.GetEntity()
	return userspostgresql.UserModel{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: gorm.DeletedAt{Time: time.Time{}, Valid: false},
	}
}

// GetUserModelWithDeletedAt returns a PostgreSQL user model with a DeletedAt set.
func GetUserModelWithDeletedAt() userspostgresql.UserModel {
	model := GetUserModel()
	entity := testusers.GetEntityWithDeletedAt()
	model.DeletedAt = gorm.DeletedAt{Time: entity.DeletedAt, Valid: true}
	return model
}

// GetWalletDB returns a PostgreSQL wallet DB.
func GetWalletDB() walletspostgresql.WalletDB {
	return walletspostgresql.NewWalletDB(GetDB())
}

// GetMockedWalletDB returns a mocked PostgreSQL user DB.
func GetMockedWalletDB() (walletspostgresql.WalletDB, sqlmock.Sqlmock, error) {
	db, mock, err := GetMockedDB()
	if err != nil {
		return walletspostgresql.WalletDB{}, nil, err
	}
	return walletspostgresql.NewWalletDB(db), mock, nil
}

// GetWalletModel returns a wallet model.
func GetWalletModel() walletspostgresql.WalletModel {
	entity := tests.GetWallet()
	model := walletspostgresql.WalletModel{
		ID:        entity.ID,
		UserID:    entity.UserID,
		Balance:   entity.Balance,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: gorm.DeletedAt{Time: time.Time{}, Valid: false},
	}
	return model
}

// GetWalletModelWithDeletedAt returns a wallet model with DeleteAt set.
func GetWalletModelWithDeletedAt() walletspostgresql.WalletModel {
	model := GetWalletModel()
	entity := tests.GetWalletWithDeletedAt()
	model.DeletedAt = gorm.DeletedAt{Time: entity.DeletedAt, Valid: true}
	return model
}

// CheckWalletsModelEntity compare if a wallet model and a wallet entity are equals
func CheckWalletsModelEntity(model walletspostgresql.WalletModel, entity wallets.Wallet) error {
	if model.ID != entity.ID {
		return fmt.Errorf("model.ID is %v, expected %v", model.ID, entity.ID)
	}
	if model.UserID != entity.UserID {
		return fmt.Errorf("model.UserID is %v, expected %v", model.UserID, entity.UserID)
	}
	if entity.Balance != model.Balance {
		return fmt.Errorf("model.Balance is %v, expected %v", model.Balance, entity.Balance)
	}
	if model.CreatedAt != entity.CreatedAt {
		return fmt.Errorf("model.CreatedAt is %v, expected %v", model.CreatedAt, entity.CreatedAt)
	}
	if model.UpdatedAt != entity.UpdatedAt {
		return fmt.Errorf("model.UpdatedAt %v, expected %v", model.UpdatedAt, entity.UpdatedAt)
	}
	if model.DeletedAt.Time != entity.DeletedAt {
		return fmt.Errorf("model.DeletedAt is %v, expected %v", model.DeletedAt.Time, entity.DeletedAt)
	}
	return nil
}

// CheckWalletsEntityModel compare if a wallet entity and wallet model are equals
func CheckWalletsEntityModel(entity wallets.Wallet, model walletspostgresql.WalletModel) error {
	if entity.ID != model.ID {
		return fmt.Errorf("wallet.ID is %v, expected %v", entity.ID, model.ID)
	}
	if entity.UserID != model.UserID {
		return fmt.Errorf("wallet.UserID is %v, expected %v", entity.UserID, model.UserID)
	}
	if entity.Balance != model.Balance {
		return fmt.Errorf("wallet.Balance is %v, expected %v", entity.Balance, model.Balance)
	}
	if entity.CreatedAt != model.CreatedAt {
		return fmt.Errorf("user.CreatedAt is %v, expected %v", entity.CreatedAt, model.CreatedAt)
	}
	if entity.UpdatedAt != model.UpdatedAt {
		return fmt.Errorf("user.UpdatedAt %v, expected %v", entity.UpdatedAt, model.UpdatedAt)
	}
	if entity.DeletedAt != model.DeletedAt.Time {
		return fmt.Errorf("user.DeletedAt is %v, expected %v", entity.DeletedAt, model.DeletedAt)
	}
	return nil
}
