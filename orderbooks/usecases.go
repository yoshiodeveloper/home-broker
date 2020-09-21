package orderbooks

import (
	"home-broker/core"
	"home-broker/orders"
)

// OrderBookUseCases represents the order use cases.
type OrderBookUseCases struct {
	orderBook *OrderBook
}

// NewOrderBookUseCases returns a new OrderBookUseCases.
func NewOrderBookUseCases(orderBook *OrderBook) OrderBookUseCases {
	return OrderBookUseCases{orderBook: orderBook}
}

// WebhookResponse is the Webhook response.
type WebhookResponse struct {
	BuyOrdersCount  int64 `json:"buy_orders_count"`
	SellOrdersCount int64 `json:"sell_orders_count"`
}

// Webhook process orders updates.
// This is a non-concurrency process.
func (orderBookUC OrderBookUseCases) Webhook(orderUpdate OrderUpdate) (WebhookResponse, error) {
	order := orderUpdate.Order
	if order.AssetID == "" {
		return WebhookResponse{}, core.NewErrValidation("Order.AssetID is invalid")
	}
	if order.ExternalID == "" {
		return WebhookResponse{}, core.NewErrValidation("Order.ExternalID is invalid")
	}
	if order.ExternalTimestamp.IsZero() {
		return WebhookResponse{}, core.NewErrValidation("Order.ExternalTimestamp is invalid")
	}
	if order.Price == 0 {
		return WebhookResponse{}, core.NewErrValidation("Order.Price is invalid")
	}
	if (order.Type != orders.OrderTypeBuy) && (order.Type != orders.OrderTypeSell) {
		return WebhookResponse{}, core.NewErrValidation("Order.Type is invalid")
	}
	orderBookUC.orderBook.Lock()
	defer orderBookUC.orderBook.Unlock()

	switch orderUpdate.Type {
	case "added":
		orderBookUC.orderBook.AddOrder(orderUpdate.Order)
	case "deleted":
		orderBookUC.orderBook.RemoveOrder(orderUpdate.Order)
	case "traded":
		orderBookUC.orderBook.DecOrderAmount(orderUpdate.Order)
	}
	return WebhookResponse{
		BuyOrdersCount:  orderBookUC.orderBook.OrdersCount[orders.OrderTypeBuy],
		SellOrdersCount: orderBookUC.orderBook.OrdersCount[orders.OrderTypeSell],
	}, nil
}
