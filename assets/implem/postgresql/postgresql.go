package postgresql

import (
	"errors"
	"home-broker/assets"
	"home-broker/core/implem/postgresql"
	"time"

	"gorm.io/gorm"
)

// AssetModel is the ORM version of Asset entity.
type AssetModel struct {
	gorm.Model
	ID         assets.AssetID    `gorm:"primaryKey;autoIncrement:true"`
	Name       string            `gorm:"not null"`
	ExchangeID assets.ExchangeID `gorm:"not null;index"`
	CreatedAt  time.Time         `gorm:"not null;index:,sort:desc"`
	UpdatedAt  time.Time         `gorm:"not null;index:,sort:desc"`
	DeletedAt  gorm.DeletedAt    `gorm:"index:,sort:desc"`
}

// TableName returns the real table name of Asset.
// It is used by GORM to perfom operations on wallet table (queries, migrations, etc.).
func (AssetModel) TableName() string {
	return "asset"
}

// AssetDB handles database commands for wallet table.
type AssetDB struct {
	assets.AssetDBInterface
	db postgresql.DB
}

// NewAssetDB creates a new AssetDB.
func NewAssetDB(db postgresql.DB) AssetDB {
	return AssetDB{db: db}
}

// ToEntity returns a Asset entity from the ORM model.
func (AssetDB) ToEntity(model AssetModel) assets.Asset {
	// "model.DeletedAt" is not a Time object. It is a struct with Time and Valid fields.
	deletedAt := time.Time{} // A "time.Time" with zero value represents a "null".
	if model.DeletedAt.Valid {
		// "model.DeletedAt" is not a "null" value.
		deletedAt = model.DeletedAt.Time
	}
	entity := assets.Asset{
		ID:         model.ID,
		Name:       model.Name,
		ExchangeID: model.ExchangeID,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
		DeletedAt:  deletedAt,
	}
	return entity
}

// ToModel returns a GORM model from an asset entity.
func (AssetDB) ToModel(entity assets.Asset) AssetModel {
	deletedAt := gorm.DeletedAt{Time: entity.DeletedAt}
	if !entity.DeletedAt.IsZero() {
		deletedAt.Valid = true
	}
	model := AssetModel{
		ID:         entity.ID,
		Name:       entity.Name,
		ExchangeID: entity.ExchangeID,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
		DeletedAt:  deletedAt,
	}
	return model
}

// GetByID returns an asset from the database by ID.
// A nil entity will be returned if it does not exist.
func (assetDB AssetDB) GetByID(id assets.AssetID) (*assets.Asset, error) {
	model := AssetModel{}
	res := assetDB.db.GetDB().Take(&model, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	entity := assetDB.ToEntity(model)
	return &entity, nil
}

// Insert inserts a new asset.
// A nil entity will be returned if an error occurs.
func (assetDB AssetDB) Insert(entity assets.Asset) (*assets.Asset, error) {
	model := assetDB.ToModel(entity)
	res := assetDB.db.GetDB().Create(&model)
	if res.Error != nil {
		return nil, res.Error
	}
	newEntity := assetDB.ToEntity(model)
	return &newEntity, nil
}
