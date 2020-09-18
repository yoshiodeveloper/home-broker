package domain

import (
	"time"
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
	UserID            UserID          // Internal user ID.
	AssetID           AssetID         // Internal asset ID.
	ExternalID        ExternalOrderID // External order ID (from a Stock Exchange)
	ExternalTimestamp time.Time       // External timestamp.
	Amount            int64
	Price             Money
	Type              OrderType
	Status            OrderStatus
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         time.Time
}

// NewBuyOrder creates a new buying order.
func NewBuyOrder(assetID AssetID, amount int64, price Money) Order {
	return Order{
		AssetID: assetID,
		Type:    OrderTypeBuy,
		Amount:  amount,
		Price:   price,
	}
}

// NewSellOrder creates a new selling order.
func NewSellOrder(assetID AssetID, amount int64, price Money) Order {
	return Order{
		AssetID: assetID,
		Type:    OrderTypeSell,
		Amount:  amount,
		Price:   price,
	}
}

// BetterThanOnSamePriceLevel returns true if this order is better offer than order parameter.
// This is only valid inside a price level (same price).
func (o *Order) BetterThanOnSamePriceLevel(order *Order) bool {
	if o.ExternalTimestamp.Before(order.ExternalTimestamp) {
		return true
	}
	return false
}

// PriceLevelOrder is a struct for an order inside OrderBookPriceLevel.
type PriceLevelOrder struct {
	Left  *PriceLevelOrder
	Right *PriceLevelOrder
	Order Order
}

// PriceLevel is a struct for a price level of buying or selling orders.
//   A price level is a group of orders (PriceLevelOrder) with the same price.
// Ex:
//   $99.00 -> PriceLevel[order, order, ...]
//   $98.00 -> PriceLevel[order, order, ...]
type PriceLevel struct {
	Left      *PriceLevel
	Right     *PriceLevel
	OrderHead *PriceLevelOrder // linked list
	Price     Money
	Type      OrderType
	// Total sum of amount of this price level.
	AmountSum int64
	// Total count of orders of this price level.
	OrdersCount int64
}

// BetterOfferThan returns true if it is a better offer than priceLevel parameter.
// Both orders must be the same Type.
func (pl *PriceLevel) BetterOfferThan(priceLevel *PriceLevel) bool {
	if pl.Type == OrderTypeBuy && pl.Price > priceLevel.Price {
		return true
	}
	if pl.Type == OrderTypeSell && pl.Price < priceLevel.Price {
		return true
	}
	return false
}

// OrderBook holds the buying and selling orders of an asset.
type OrderBook struct {
	AssetID AssetID

	// ordersByOrderID maps an Order ID to a PriceLevelOrder.
	// This improves the search of a order inside the order book at O(1).
	// The downside is the rehash of a hashmap.
	// This can be improved with a AVL BTree.
	OrdersByOrderID map[OrderID]*PriceLevelOrder

	PriceLevelsByPrices map[OrderType]map[Money]*PriceLevel
	PriceLevelsHeads    map[OrderType]*PriceLevel // linked list for buying and selling
	OrdersCount         map[OrderType]int64
}

// NewOrderBook creates a new OrderBook.
func NewOrderBook(assetID AssetID) *OrderBook {
	ob := OrderBook{
		AssetID:         assetID,
		OrdersByOrderID: make(map[OrderID]*PriceLevelOrder),
		PriceLevelsByPrices: map[OrderType]map[Money]*PriceLevel{
			OrderTypeBuy:  make(map[Money]*PriceLevel),
			OrderTypeSell: make(map[Money]*PriceLevel),
		},
		PriceLevelsHeads: make(map[OrderType]*PriceLevel),
		OrdersCount:      make(map[OrderType]int64),
	}
	return &ob
}

// addNewPriceLevel adds a new PriceLevel into the OrderBook.
func (ob *OrderBook) addNewPriceLevel(order *Order) *PriceLevel {
	priceLevel := &PriceLevel{Price: order.Price, Type: order.Type}
	ob.PriceLevelsByPrices[priceLevel.Type][priceLevel.Price] = priceLevel

	if ob.PriceLevelsHeads[order.Type] == nil {
		ob.PriceLevelsHeads[order.Type] = priceLevel // head
		return priceLevel
	}

	currPL := ob.PriceLevelsHeads[order.Type] // head
	for {
		if priceLevel.BetterOfferThan(currPL) {
			if currPL.Left == nil {
				// currPL is the first node.
				ob.PriceLevelsHeads[order.Type] = priceLevel
			} else {
				currPL.Left.Right = priceLevel
			}
			priceLevel.Left = currPL.Left
			priceLevel.Right = currPL
			currPL.Left = priceLevel
			break
		}
		if currPL.Right == nil {
			// currPL is the last node.
			currPL.Right = priceLevel
			priceLevel.Left = currPL
			break
		}
		currPL = currPL.Right
	}
	return priceLevel
}

func (ob *OrderBook) addNewPriceLevelOrder(order *Order) {
	priceLevel, _ := ob.PriceLevelsByPrices[order.Type][order.Price]
	if priceLevel == nil {
		priceLevel = ob.addNewPriceLevel(order)
	}

	priceLevel.AmountSum += order.Amount
	priceLevel.OrdersCount++

	plOrder := &PriceLevelOrder{Order: *order}
	ob.OrdersByOrderID[plOrder.Order.ID] = plOrder
	ob.OrdersCount[order.Type]++

	if priceLevel.OrderHead == nil {
		priceLevel.OrderHead = plOrder
		return
	}

	currPLOrder := priceLevel.OrderHead
	for {
		if order.BetterThanOnSamePriceLevel(&currPLOrder.Order) {
			if currPLOrder.Left == nil {
				// currPLOrder is the first node.
				priceLevel.OrderHead = plOrder
			} else {
				currPLOrder.Left.Right = plOrder
			}
			plOrder.Left, plOrder.Right = currPLOrder.Left, currPLOrder
			currPLOrder.Left = plOrder
			break
		}
		if currPLOrder.Right == nil {
			// currPLOrder is the last node.
			currPLOrder.Right, plOrder.Left = plOrder, currPLOrder
			break
		}
		currPLOrder = currPLOrder.Right
	}
}

// AddOrder adds a order into the OrderBook.
func (ob *OrderBook) AddOrder(order Order) {
	// We receive the order as a copy we can pass on as reference.
	ob.addNewPriceLevelOrder(&order)
}

// GetBuyOrders returns a slice of buying orders.
func (ob *OrderBook) GetBuyOrders() []Order {
	orders := ob.getOrders(OrderTypeBuy)
	return *orders
}

// GetSellOrders returns a slice of selling orders.
func (ob *OrderBook) GetSellOrders() []Order {
	orders := ob.getOrders(OrderTypeSell)
	return *orders
}

// getOrders returns a slice of buying or selling orders.
func (ob *OrderBook) getOrders(orderType OrderType) *[]Order {
	currPriceLevel := ob.PriceLevelsHeads[orderType]
	orders := make([]Order, 0, ob.OrdersCount[orderType])
	for currPriceLevel != nil {
		currPLOrder := currPriceLevel.OrderHead
		for currPLOrder != nil {
			orders = append(orders, currPLOrder.Order)
			currPLOrder = currPLOrder.Right
		}
		currPriceLevel = currPriceLevel.Right
	}
	return &orders
}

// RemoveOrder removes an order into the order book.
// func (ob *OrderBook) RemoveOrder(orderID OrderID)
