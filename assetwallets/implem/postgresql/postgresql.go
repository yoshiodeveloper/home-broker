package postgresql

import (
	"errors"
	"home-broker/assets"
	assetspostgresql "home-broker/assets/implem/postgresql"
	"home-broker/assetwallets"
	"home-broker/core/implem/postgresql"
	"home-broker/users"
	userspostgresql "home-broker/users/implem/postgresql"
	"strings"
	"time"

	"gorm.io/gorm"
)

// AssetWalletModel is the ORM version of AssetWallet entity.
type AssetWalletModel struct {
	gorm.Model
	ID        assetwallets.AssetWalletID `gorm:"primaryKey;autoIncrement:true"`
	UserID    users.UserID               `gorm:"uniqueIndex:idx_assetwallet_userasset;not null"`
	User      userspostgresql.UserModel
	AssetID   assets.AssetID `gorm:"uniqueIndex:idx_assetwallet_userasset;not null"`
	Asset     assetspostgresql.AssetModel
	Balance   assets.AssetUnit `gorm:"not null;index:,sort:desc"` // Mind the money.MoneyDecimalPlaces.
	CreatedAt time.Time        `gorm:"not null;index:,sort:desc"`
	UpdatedAt time.Time        `gorm:"not null;index:,sort:desc"`
	DeletedAt gorm.DeletedAt   `gorm:"index:,sort:desc"`
}

// TableName returns the real table name of AssetWallet.
// It is used by GORM to perfom operations on asset wallet table (queries, migrations, etc.).
func (AssetWalletModel) TableName() string {
	return "walletasset"
}

// AssetWalletDB handles database commands for asset wallet table.
type AssetWalletDB struct {
	assetwallets.AssetWalletDBInterface
	db postgresql.DB
}

// NewAssetWalletDB creates a new AssetWalletDB.
func NewAssetWalletDB(db postgresql.DB) AssetWalletDB {
	return AssetWalletDB{db: db}
}

// ToEntity returns a Wallet entity from the ORM model.
func (AssetWalletDB) ToEntity(model AssetWalletModel) assetwallets.AssetWallet {
	// "model.DeletedAt" is not a Time object. It is a struct with Time and Valid fields.
	deletedAt := time.Time{} // A "time.Time" with zero value represents a "null".
	if model.DeletedAt.Valid {
		// "model.DeletedAt" is not a "null" value.
		deletedAt = model.DeletedAt.Time
	}
	entity := assetwallets.AssetWallet{
		ID:        model.ID,
		UserID:    model.UserID,
		AssetID:   model.AssetID,
		Balance:   model.Balance,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: deletedAt,
	}
	return entity
}

// ToModel returns a GORM model from a asset wallet entity.
func (AssetWalletDB) ToModel(entity assetwallets.AssetWallet) AssetWalletModel {
	deletedAt := gorm.DeletedAt{Time: entity.DeletedAt}
	if !entity.DeletedAt.IsZero() {
		deletedAt.Valid = true
	}
	model := AssetWalletModel{
		ID:        entity.ID,
		UserID:    entity.UserID,
		AssetID:   entity.AssetID,
		Balance:   entity.Balance,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: deletedAt,
	}
	return model
}

// GetByUserIDAssetID returns an asset wallet from the database by an user ID and asset ID.
// A nil entity will be returned if it does not exist.
func (assetWalletDB AssetWalletDB) GetByUserIDAssetID(userID users.UserID, assetID assets.AssetID) (*assetwallets.AssetWallet, error) {
	model := AssetWalletModel{}
	modelCond := AssetWalletModel{UserID: userID, AssetID: assetID}
	res := assetWalletDB.db.GetDB().Take(&model, modelCond)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	entity := assetWalletDB.ToEntity(model)
	return &entity, nil
}

// Insert inserts a new asset wallet.
// A nil entity will be returned if an error occurs.
// The following errors can happen: ErrAssetWalletAlreadyExists, ErrAssetDoesNotExist, ErrUserDoesNotExist.
func (assetWalletDB AssetWalletDB) Insert(entity assetwallets.AssetWallet) (*assetwallets.AssetWallet, error) {
	model := assetWalletDB.ToModel(entity)
	res := assetWalletDB.db.GetDB().Create(&model)
	if res.Error != nil {
		errMsg := res.Error.Error()
		if strings.Contains(errMsg, "unique constraint") {
			// Original error: "ERROR: duplicate key value violates unique constraint "assetwallet_pkey" (SQLSTATE 23505)"
			// Original error: "ERROR: duplicate key value violates unique constraint "idx_assetwallet_userasset" (SQLSTATE 23505)"
			return nil, assetwallets.ErrAssetWalletAlreadyExists
		}
		if strings.Contains(errMsg, "foreign key constraint") && strings.Contains(errMsg, "user") {
			// Original error: "ERROR: insert or update on table "assetwallet" violates foreign key constraint "fk_assetwallet_user" (SQLSTATE 23503)"
			return nil, users.ErrUserDoesNotExist
		}
		if strings.Contains(errMsg, "foreign key constraint") && strings.Contains(errMsg, "asset") {
			// Original error: "ERROR: insert or update on table "assetwallet" violates foreign key constraint "fk_assetwallet_asset" (SQLSTATE 23503)"
			return nil, assets.ErrAssetDoesNotExist
		}
		return nil, res.Error
	}
	newEntity := assetWalletDB.ToEntity(model)
	return &newEntity, nil
}

// IncBalanceByUserIDAssetID increments or decrements the balance field by a user ID and asset ID.
// An updated entity entity will be returned or nil if it does not exist.
func (assetWalletDB AssetWalletDB) IncBalanceByUserIDAssetID(userID users.UserID, assetID assets.AssetID, amount assets.AssetUnit) (*assetwallets.AssetWallet, error) {
	updatedAt := time.Now()
	res := assetWalletDB.db.GetDB().
		Table("assetwallet").
		Where(`"user_id"=? AND "asset_id"=? "deleted_at" IS NULL`, userID, assetID).
		Updates(map[string]interface{}{
			"balance":    gorm.Expr(`"balance"+?`, amount),
			"updated_at": updatedAt,
		})
	if res.Error != nil {
		return nil, res.Error
	}
	entity, err := assetWalletDB.GetByUserIDAssetID(userID, assetID)
	return entity, err
}
