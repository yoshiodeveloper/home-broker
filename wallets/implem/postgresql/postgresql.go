package postgresql

import (
	"errors"
	"fmt"
	"home-broker/core/implem/postgresql"
	"home-broker/money"
	"home-broker/users"
	userspostgresql "home-broker/users/implem/postgresql"
	"home-broker/wallets"
	"strings"
	"time"

	"gorm.io/gorm"
)

// WalletModel is the ORM version of Wallet entity.
type WalletModel struct {
	gorm.Model
	ID        wallets.WalletID `gorm:"primaryKey;autoIncrement:true"`
	UserID    users.UserID     `gorm:"unique;not null"`
	User      userspostgresql.UserModel
	Balance   money.Money    `gorm:"not null;index:,sort:desc"` // Mind the money.MoneyDecimalPlaces.
	CreatedAt time.Time      `gorm:"not null;index:,sort:desc"`
	UpdatedAt time.Time      `gorm:"not null;index:,sort:desc"`
	DeletedAt gorm.DeletedAt `gorm:"index:,sort:desc"`
}

// TableName returns the real table name of Wallet.
// It is used by GORM to perfom operations on wallet table (queries, migrations, etc.).
func (WalletModel) TableName() string {
	return "wallet"
}

// WalletDB handles database commands for wallet table.
type WalletDB struct {
	wallets.WalletDBInterface
	db postgresql.DB
}

// NewWalletDB creates a new WalletDB.
func NewWalletDB(db postgresql.DB) WalletDB {
	return WalletDB{db: db}
}

// ToEntity returns a Wallet entity from the ORM model.
func (WalletDB) ToEntity(model WalletModel) wallets.Wallet {
	// "model.DeletedAt" is not a Time object. It is a struct with Time and Valid fields.
	deletedAt := time.Time{} // A "time.Time" with zero value represents a "null".
	if model.DeletedAt.Valid {
		// "model.DeletedAt" is not a "null" value.
		deletedAt = model.DeletedAt.Time
	}
	entity := wallets.Wallet{
		ID:        model.ID,
		UserID:    model.UserID,
		Balance:   model.Balance,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: deletedAt,
	}
	return entity
}

// ToModel returns a GORM model from an wallet entity.
func (WalletDB) ToModel(entity wallets.Wallet) WalletModel {
	deletedAt := gorm.DeletedAt{Time: entity.DeletedAt}
	if !entity.DeletedAt.IsZero() {
		deletedAt.Valid = true
	}
	model := WalletModel{
		ID:        entity.ID,
		UserID:    entity.UserID,
		Balance:   entity.Balance,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: deletedAt,
	}
	return model
}

// GetByUserID returns a wallet from the database by an user ID.
// A nil entity will be returned if it does not exist.
func (walletDB WalletDB) GetByUserID(userID users.UserID) (*wallets.Wallet, error) {
	model := WalletModel{}
	modelCond := WalletModel{UserID: userID}
	res := walletDB.db.GetDB().Take(&model, modelCond)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	entity := walletDB.ToEntity(model)
	return &entity, nil
}

// Insert inserts a new wallet.
// A nil entity will be returned if an error occurs.
// The following errors can happen: ErrWalletAlreadyExists, ErrUserDoesNotExist.
func (walletDB WalletDB) Insert(entity wallets.Wallet) (*wallets.Wallet, error) {
	model := walletDB.ToModel(entity)
	res := walletDB.db.GetDB().Create(&model)
	if res.Error != nil {
		errMsg := res.Error.Error()
		if strings.Contains(errMsg, "unique constraint") && entity.ID != 0 {
			// Original error: "ERROR: duplicate key value violates unique constraint "wallet_pkey" (SQLSTATE 23505)"
			return nil, fmt.Errorf("%w: ID %d", wallets.ErrWalletAlreadyExists, entity.ID)
		}
		if strings.Contains(errMsg, "unique constraint") && strings.Contains(errMsg, "user_id") {
			// Original error: "ERROR: duplicate key value violates unique constraint "wallet_user_id_key" (SQLSTATE 23505)"
			return nil, fmt.Errorf("%w: User ID %d", wallets.ErrWalletAlreadyExists, entity.UserID)
		}
		if strings.Contains(errMsg, "foreign key constraint") && strings.Contains(errMsg, "user") {
			// Original error: "ERROR: insert or update on table "wallet" violates foreign key constraint "fk_wallet_user" (SQLSTATE 23503)"
			return nil, fmt.Errorf("%w: User ID %d", users.ErrUserDoesNotExist, entity.UserID)
		}
		return nil, res.Error
	}
	newEntity := walletDB.ToEntity(model)
	return &newEntity, nil
}

// IncBalanceByUserID increments or decrements the balance field by a user ID.
// An updated entity entity will be returned or nil if it does not exist.
func (walletDB WalletDB) IncBalanceByUserID(userID users.UserID, amount money.Money) (*wallets.Wallet, error) {
	updatedAt := time.Now()
	res := walletDB.db.GetDB().
		Table("wallet").
		Where(`"user_id"=? AND "deleted_at" IS NULL`, userID).
		Updates(map[string]interface{}{
			"balance":    gorm.Expr(`"balance"+?`, amount),
			"updated_at": updatedAt,
		})
	if res.Error != nil {
		return nil, res.Error
	}
	entity, err := walletDB.GetByUserID(userID)
	return entity, err
}
