package postgresql

import (
	"errors"
	"home-broker/assets"
	assetspostgresql "home-broker/assets/implem/postgresql"
	"home-broker/core/implem/postgresql"
	"home-broker/money"
	"home-broker/orders"
	"home-broker/users"
	userspostgresql "home-broker/users/implem/postgresql"
	"strings"
	"time"

	"gorm.io/gorm"
)

// OrderModel is the ORM version of Order entity.
type OrderModel struct {
	gorm.Model
	ID                orders.OrderID `gorm:"primaryKey;autoIncrement:true"`
	UserID            users.UserID   `gorm:"not null;index:,sort:desc"`
	User              userspostgresql.UserModel
	AssetID           assets.AssetID `gorm:"not null;index:,sort:desc"`
	Asset             assetspostgresql.AssetModel
	ExternalID        orders.ExternalOrderID `gorm:"not null;index"`
	ExternalTimestamp time.Time              `gorm:"not null;index:,sort:desc"`
	Amount            assets.AssetUnit       `gorm:"not null"`
	Price             money.Money            `gorm:"not null"`
	Type              orders.OrderType       `gorm:"not null;index"`
	Status            orders.OrderStatus     `gorm:"not null"`
	CreatedAt         time.Time              `gorm:"not null;index:,sort:desc"`
	UpdatedAt         time.Time              `gorm:"not null;index:,sort:desc"`
	DeletedAt         gorm.DeletedAt         `gorm:"index:,sort:desc"`
}

// TableName returns the real table name of Order.
// It is used by GORM to perfom operations on wallet table (queries, migrations, etc.).
func (OrderModel) TableName() string {
	return "order"
}

// OrderDB handles database commands for wallet table.
type OrderDB struct {
	orders.OrderDBInterface
	db postgresql.DB
}

// NewOrderDB creates a new OrderDB.
func NewOrderDB(db postgresql.DB) OrderDB {
	return OrderDB{db: db}
}

// ToEntity returns a Order entity from the ORM model.
func (OrderDB) ToEntity(model OrderModel) orders.Order {
	// "model.DeletedAt" is not a Time object. It is a struct with Time and Valid fields.
	deletedAt := time.Time{} // A "time.Time" with zero value represents a "null".
	if model.DeletedAt.Valid {
		// "model.DeletedAt" is not a "null" value.
		deletedAt = model.DeletedAt.Time
	}
	entity := orders.Order{
		ID:                model.ID,
		UserID:            model.UserID,
		AssetID:           model.AssetID,
		ExternalID:        model.ExternalID,
		ExternalTimestamp: model.ExternalTimestamp,
		Amount:            model.Amount,
		Price:             model.Price,
		Type:              model.Type,
		Status:            model.Status,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
		DeletedAt:         deletedAt,
	}
	return entity
}

// ToModel returns a GORM model from an order entity.
func (OrderDB) ToModel(entity orders.Order) OrderModel {
	deletedAt := gorm.DeletedAt{Time: entity.DeletedAt}
	if !entity.DeletedAt.IsZero() {
		deletedAt.Valid = true
	}
	model := OrderModel{
		ID:                entity.ID,
		UserID:            entity.UserID,
		AssetID:           entity.AssetID,
		ExternalID:        entity.ExternalID,
		ExternalTimestamp: entity.ExternalTimestamp,
		Amount:            entity.Amount,
		Price:             entity.Price,
		Type:              entity.Type,
		Status:            entity.Status,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
		DeletedAt:         deletedAt,
	}
	return model
}

// GetByID returns an order by ID.
// A nil entity will be returned if it does not exist.
func (orderDB OrderDB) GetByID(id orders.OrderID) (*orders.Order, error) {
	model := OrderModel{}
	res := orderDB.db.GetDB().Take(&model, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	entity := orderDB.ToEntity(model)
	return &entity, nil
}

// GetByExternalIDAssetID returns an order by external ID and asset ID.
// If the record does not exist a nil entity will be returned.
func (orderDB OrderDB) GetByExternalIDAssetID(externalID orders.ExternalOrderID, assetID assets.AssetID) (*orders.Order, error) {
	model := OrderModel{}
	res := orderDB.db.GetDB().Where(`"external_id"=? AND "asset_id"=?`, externalID, assetID).Take(&model)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	entity := orderDB.ToEntity(model)
	return &entity, nil
}

// Insert inserts a new order.
// A nil entity will be returned if an error occurs.
// The following errors can happen: ErrUserDoesNotExist.
func (orderDB OrderDB) Insert(entity orders.Order) (*orders.Order, error) {
	model := orderDB.ToModel(entity)
	res := orderDB.db.GetDB().Create(&model)
	if res.Error != nil {
		errMsg := res.Error.Error()
		if strings.Contains(errMsg, "foreign key constraint") && strings.Contains(errMsg, "asset") {
			// Original error: "ERROR: insert or update on table "order" violates foreign key constraint "fk_order_asset" (SQLSTATE 23503)"
			return nil, assets.ErrAssetDoesNotExist
		}
		if strings.Contains(errMsg, "foreign key constraint") && strings.Contains(errMsg, "user") {
			// Original error: "ERROR: insert or update on table "order" violates foreign key constraint "fk_order_user" (SQLSTATE 23503)"
			return nil, users.ErrUserDoesNotExist
		}
		return nil, res.Error
	}
	newEntity := orderDB.ToEntity(model)
	return &newEntity, nil
}

// UpdateExternalResponse updates an order base on a exchange response.
func (orderDB OrderDB) UpdateExternalResponse(orderID orders.OrderID, externalID orders.ExternalOrderID, externalTimestamp time.Time, status orders.OrderStatus) error {
	updatedAt := time.Now()
	res := orderDB.db.GetDB().
		Table("order").
		Where(`"id"=?`, orderID).
		Updates(map[string]interface{}{
			"external_id":        externalID,
			"external_timestamp": externalTimestamp,
			"status":             status,
			"updated_at":         updatedAt,
		})
	return res.Error
}

// UpdateStatus updates an order status.
func (orderDB OrderDB) UpdateStatus(orderID orders.OrderID, status orders.OrderStatus) error {
	updatedAt := time.Now()
	res := orderDB.db.GetDB().
		Table("order").
		Where(`"id"=?`, orderID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": updatedAt,
		})
	return res.Error
}
