package postgresql

import (
	"errors"
	"fmt"
	"home-broker/domain"
	"home-broker/infra"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// WalletGORM is the ORM version of Wallet entity.
type WalletGORM struct {
	gorm.Model
	ID        int64 `gorm:"primaryKey;autoIncrement:true"`
	UserID    int64 `gorm:"unique;not null"`
	User      UserGORM
	Balance   decimal.Decimal `gorm:"not null;type:NUMERIC(21,6);index:,sort:desc"` // Mind the domain.MoneyDecimalPlaces.
	CreatedAt time.Time       `gorm:"not null;index:,sort:desc"`
	UpdatedAt time.Time       `gorm:"not null;index:,sort:desc"`
	DeletedAt gorm.DeletedAt  `gorm:"index:,sort:desc"`
}

// TableName returns the real table name of Wallet.
// It is used by GORM to perfom operations on wallet table (queries, migrations, etc.).
func (WalletGORM) TableName() string {
	return "wallet"
}

// WalletRepo handles database commands for wallet table.
type WalletRepo struct {
	infra.WalletRepoI
	dbClient *DBClient
}

// NewWalletRepo creates a new WalletRepo.
func NewWalletRepo(dbClient *DBClient) *WalletRepo {
	return &WalletRepo{
		dbClient: dbClient,
	}
}

// ToEntity returns a Wallet entity from the ORM model.
func (repo *WalletRepo) ToEntity(modelGORM *WalletGORM) domain.Wallet {
	// "modelGORM.DeletedAt" is not a Time object. It is a struct with Time and Valid fields.
	deletedAt := time.Time{} // A "time.Time" with zero value represents a "null".
	if modelGORM.DeletedAt.Valid {
		// "model.DeletedAt" is not a "null" value.
		deletedAt = modelGORM.DeletedAt.Time
	}
	entity := domain.Wallet{
		ID:        domain.WalletID(modelGORM.ID),
		UserID:    domain.UserID(modelGORM.UserID),
		Balance:   domain.Money{Number: modelGORM.Balance},
		CreatedAt: modelGORM.CreatedAt,
		UpdatedAt: modelGORM.UpdatedAt,
		DeletedAt: deletedAt,
	}
	return entity
}

// ToGORMModel returns a GORM model from an wallet entity.
func (repo *WalletRepo) ToGORMModel(entity *domain.Wallet) WalletGORM {
	deletedAt := gorm.DeletedAt{Time: entity.DeletedAt}
	if !entity.DeletedAt.IsZero() {
		deletedAt.Valid = true
	}
	model := WalletGORM{
		ID:        int64(entity.ID),
		UserID:    int64(entity.UserID),
		Balance:   entity.Balance.Number,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: deletedAt,
	}
	return model
}

// GetByUserID returns a wallet from the database by an user ID.
func (repo *WalletRepo) GetByUserID(id domain.UserID) (*domain.Wallet, error) {
	modelGORM := WalletGORM{}
	modelGORMCond := WalletGORM{UserID: int64(id)}
	res := repo.dbClient.db.Take(&modelGORM, modelGORMCond)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	entity := repo.ToEntity(&modelGORM)
	return &entity, nil
}

// Insert inserts a new wallet.
func (repo *WalletRepo) Insert(entity domain.Wallet) (*domain.Wallet, error) {
	modelGORM := repo.ToGORMModel(&entity)
	res := repo.dbClient.db.Create(&modelGORM)
	if res.Error != nil {
		errMsg := res.Error.Error()
		if strings.Contains(errMsg, "unique constraint") && entity.ID != 0 {
			// Original error: "ERROR: duplicate key value violates unique constraint "wallet_pkey" (SQLSTATE 23505)"
			return nil, fmt.Errorf("%w: ID %d", infra.ErrWalletAlreadyExists, entity.ID)
		}
		if strings.Contains(errMsg, "unique constraint") && strings.Contains(errMsg, "user_id") {
			// Original error: "ERROR: duplicate key value violates unique constraint "wallet_user_id_key" (SQLSTATE 23505)"
			return nil, fmt.Errorf("%w: User ID %d", infra.ErrWalletAlreadyExists, entity.UserID)
		}
		if strings.Contains(errMsg, "foreign key constraint") && strings.Contains(errMsg, "user") {
			// Original error: "ERROR: insert or update on table "wallet" violates foreign key constraint "fk_wallet_user" (SQLSTATE 23503)"
			return nil, fmt.Errorf("%w: User ID %d", infra.ErrUserDoesNotExist, entity.UserID)
		}
		return nil, res.Error
	}
	newEntity := repo.ToEntity(&modelGORM)
	return &newEntity, nil
}
