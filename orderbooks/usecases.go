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
func (orderBookUC OrderBookUseCases) Webhook(externalUp orders.ExternalUpdate) (WebhookResponse, error) {
	if externalUp.AssetID == "" {
		return WebhookResponse{}, core.NewErrValidation("Asset ID is invalid")
	}
	if externalUp.ID == "" {
		return WebhookResponse{}, core.NewErrValidation("ID (external) is invalid")
	}
	if externalUp.Timestamp.IsZero() {
		return WebhookResponse{}, core.NewErrValidation("Timestamp is invalid")
	}
	if (externalUp.Type != orders.OrderTypeBuy) && (externalUp.Type != orders.OrderTypeSell) {
		return WebhookResponse{}, core.NewErrValidation("Type is invalid")
	}

	order := Order{ // this is not the same as "orders.Order" type.
		Mine:      externalUp.Mine,
		ID:        externalUp.ID,
		AssetID:   externalUp.AssetID,
		Price:     externalUp.Price,
		Amount:    externalUp.Amount,
		Type:      externalUp.Type,
		Timestamp: externalUp.Timestamp,
	}

	var tradeRequest *TradeRequest

	func() {
		orderBookUC.orderBook.Lock()
		defer orderBookUC.orderBook.Unlock()

		switch externalUp.Action {
		case "added":
			tradeRequest = orderBookUC.orderBook.AddOrder(order)

		case "deleted":
			orderBookUC.orderBook.RemoveOrder(order)
		case "traded":
			orderBookUC.orderBook.DecOrderAmount(order)
		}
	}()

	if tradeRequest != nil {
		// TODO: Do request on the exchange OR on a kafka topic.
		// We also need to change the order status.
	}

	return WebhookResponse{
		BuyOrdersCount:  orderBookUC.orderBook.OrdersCount[orders.OrderTypeBuy],
		SellOrdersCount: orderBookUC.orderBook.OrdersCount[orders.OrderTypeSell],
	}, nil
}
