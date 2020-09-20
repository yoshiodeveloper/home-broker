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
	OrderStatus int8
)

const (
	// OrderTypeBuy is an order type for a buy.
	OrderTypeBuy OrderType = "buy"

	// OrderTypeSell is an order type for a sell.
	OrderTypeSell OrderType = "sell"

	// OrderStatusAccepted is an accepted order.
	// This means that the order was accepted and processed by the exchange.
	OrderStatusAccepted int8 = 1

	// OrderStatusDenied is a denied order.
	// This happens when the order is not accepted or processed by the exchange.
	// This can happen because of an exchange error or incorrect data.
	OrderStatusDenied int8 = -1

	// OrderStatusPending is an pending order.
	// This happens when the order was sent to the exchange and we are waiting to process.
	OrderStatusPending int8 = 2
)

// Order is an entity for buying or selling intentions.
type Order struct {
	ID                OrderID         // Internal ID.
	UserID            users.UserID    // Internal user ID.
	AssetID           assets.AssetID  // Internal asset ID.
	ExternalID        ExternalOrderID // External order ID (from a Stock Exchange)
	ExternalTimestamp time.Time       // External timestamp.
	Amount            int64
	Price             money.Money
	Type              OrderType
	Status            OrderStatus
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
}

// NewBuyOrder creates a new buying order.
func NewBuyOrder(assetID assets.AssetID, amount int64, price money.Money) Order {
	return Order{
		AssetID: assetID,
		Type:    OrderTypeBuy,
		Amount:  amount,
		Price:   price,
	}
}

// NewSellOrder creates a new selling order.
func NewSellOrder(assetID assets.AssetID, amount int64, price money.Money) Order {
	return Order{
		AssetID: assetID,
		Type:    OrderTypeSell,
		Amount:  amount,
		Price:   price,
	}
}

// BetterThan returns true if this order is better offer than order parameter.
func (o Order) BetterThan(order Order) bool {
	/*
		// These checks are not been used, as we check the order already inside a price level (at same price).
		if o.Type == OrderTypeBuy && o.Price > order.Price {
			return true
		}
		if o.Type == OrderTypeSell && o.Price < order.Price {
			return true
		}
	*/
	if o.ExternalTimestamp.Before(order.ExternalTimestamp) {
		return true
	}
	return false
}
