package orders

import (
	"home-broker/assets"
	"home-broker/money"
	"home-broker/users"
	"time"
)

type (
	// OrderType represents a order type.
	// Use the value of OrderTypeBuy or OrderTypeSell to set this data type.
	OrderType string

	// OrderID represents the Order ID type.
	OrderID int64

	// ExternalOrderID represents an external order ID (ex: an order ID from a exchange).
	ExternalOrderID string

	// OrderStatus represents the status of an order.
	OrderStatus string
)

const (
	// OrderTypeBuy is an order type for a buy.
	OrderTypeBuy OrderType = "buy"

	// OrderTypeSell is an order type for a sell.
	OrderTypeSell OrderType = "sell"

	// OrderStatusAccepted is an accepted order.
	// This means that the order was accepted and processed by the exchange.
	OrderStatusAccepted = "accepted"

	// OrderStatusPending is an pending order.
	// This happens when the order was sent to the exchange and we are waiting to process.
	OrderStatusPending = "pending"

	// OrderStatusDenied is a denied order.
	// This happens when the order is not accepted or processed by the exchange.
	// This can happen because of an exchange error or incorrect data.
	OrderStatusDenied = "denied"

	// OrderStatusCanceling is a canceling order.
	OrderStatusCanceling = "canceling"

	// OrderStatusCanceled is a canceled order.
	OrderStatusCanceled = "canceled"

	// ExternalUpdateActionAdded is an order added to the order book.
	ExternalUpdateActionAdded = "added"

	// ExternalUpdateActionDeleted is an order removed from the order book.
	ExternalUpdateActionDeleted = "deleted"

	// ExternalUpdateActionTraded is an order traded and must be removed from the order book.
	ExternalUpdateActionTraded = "traded"
)

// Order is an entity for buying or selling intentions.
type Order struct {
	ID                OrderID          `json:"id"`                 // Internal ID.
	UserID            users.UserID     `json:"user_id"`            // Internal user ID.
	AssetID           assets.AssetID   `json:"asset_id"`           // Internal asset ID.
	ExternalID        ExternalOrderID  `json:"external_id"`        // External order ID (from a Stock Exchange)
	ExternalTimestamp time.Time        `json:"external_timestamp"` // External timestamp.
	Amount            assets.AssetUnit `json:"amount"`
	Price             money.Money      `json:"price"`
	Type              OrderType        `json:"type"`
	Status            OrderStatus      `json:"status"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         time.Time        `json:"-"`
}

// NewBuyOrder creates a new buying order.
func NewBuyOrder(assetID assets.AssetID, amount assets.AssetUnit, price money.Money) Order {
	return Order{
		AssetID: assetID,
		Type:    OrderTypeBuy,
		Amount:  amount,
		Price:   price,
	}
}

// NewSellOrder creates a new selling order.
func NewSellOrder(assetID assets.AssetID, amount assets.AssetUnit, price money.Money) Order {
	return Order{
		AssetID: assetID,
		Type:    OrderTypeSell,
		Amount:  amount,
		Price:   price,
	}
}

// ExternalUpdate holds an order update sent by an exchange service.
type ExternalUpdate struct {
	ID        ExternalOrderID  `json:"id"`
	AssetID   assets.AssetID   `json:"asset_id"`
	Price     money.Money      `json:"price"`
	Amount    assets.AssetUnit `json:"amount"`
	Type      OrderType        `json:"type"`
	Timestamp time.Time        `json:"timestamp"`
	Action    string           `json:"action"` // added / deleted / traded
}
