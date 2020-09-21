package orderbooks

import (
	"home-broker/assets"
	"home-broker/money"
	"home-broker/orders"
	"log"
	"sync"
	"time"
)

// Order holds an order sent by an exchange service.
// This not the same as "orders.Order".
type Order struct {
	Mine      bool                   `json:"mine"`
	ID        orders.ExternalOrderID `json:"id"`
	AssetID   assets.AssetID         `json:"asset_id"`
	Price     money.Money            `json:"price"`
	Amount    assets.AssetUnit       `json:"amount"`
	Type      orders.OrderType       `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	InTrade   bool                   `json:"in_trade"`
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
	if o.Timestamp.Before(order.Timestamp) {
		return true
	}
	return false
}

// TradeRequest represents a trade request.
type TradeRequest struct {
	InterestedOrder Order
	InterestOrder   Order
	Amount          assets.AssetUnit
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
	Price     money.Money
	Type      orders.OrderType
	// Total sum of amount of this price level.
	AmountSum assets.AssetUnit
	// Total count of orders of this price level.
	OrdersCount int64
}

// BetterOfferThan returns true if it is a better offer than priceLevel parameter.
// Both orders must be the same Type.
func (pl PriceLevel) BetterOfferThan(priceLevel PriceLevel) bool {
	if pl.Type == orders.OrderTypeBuy && pl.Price > priceLevel.Price {
		return true
	}
	if pl.Type == orders.OrderTypeSell && pl.Price < priceLevel.Price {
		return true
	}
	return false
}

// OrderBook holds the buying and selling orders of an asset.
type OrderBook struct {
	// Only one goroutine can perfom operations in this Order Book at time.
	mux sync.Mutex

	AssetID assets.AssetID

	// ordersByOrderID maps an Order ID to a PriceLevelOrder.
	// The ID is the external ID (from the exchange).
	// This improves the search of a order inside the order book at O(1).
	// The downside is the rehash of a hashmap.
	// This can be improved with a AVL BTree.
	OrdersByOrderID map[orders.ExternalOrderID]*PriceLevelOrder

	PriceLevelsByPrices map[orders.OrderType]map[money.Money]*PriceLevel
	PriceLevelsHeads    map[orders.OrderType]*PriceLevel // linked list for buying and selling
	OrdersCount         map[orders.OrderType]int64
}

// NewOrderBook creates a new OrderBook.
func NewOrderBook(assetID assets.AssetID) *OrderBook {
	ob := OrderBook{
		AssetID:         assetID,
		OrdersByOrderID: make(map[orders.ExternalOrderID]*PriceLevelOrder),
		PriceLevelsByPrices: map[orders.OrderType]map[money.Money]*PriceLevel{
			orders.OrderTypeBuy:  make(map[money.Money]*PriceLevel),
			orders.OrderTypeSell: make(map[money.Money]*PriceLevel),
		},
		PriceLevelsHeads: make(map[orders.OrderType]*PriceLevel),
		OrdersCount:      make(map[orders.OrderType]int64),
	}
	return &ob
}

// Lock locks for thread safety.
// Call this before any operation.
func (ob *OrderBook) Lock() {
	ob.mux.Lock()
}

// Unlock unlocks for thread safety.
func (ob *OrderBook) Unlock() {
	ob.mux.Unlock()
}

// addNewPriceLevel adds a new PriceLevel into the OrderBook.
func (ob *OrderBook) addNewPriceLevel(order Order) *PriceLevel {
	priceLevel := &PriceLevel{Price: order.Price, Type: order.Type}
	ob.PriceLevelsByPrices[priceLevel.Type][priceLevel.Price] = priceLevel

	if ob.PriceLevelsHeads[order.Type] == nil {
		ob.PriceLevelsHeads[order.Type] = priceLevel // head
		return priceLevel
	}

	currPL := ob.PriceLevelsHeads[order.Type] // head
	for {
		if priceLevel.BetterOfferThan(*currPL) {
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

func (ob *OrderBook) addNewPriceLevelOrder(order Order) *PriceLevelOrder {
	priceLevel, _ := ob.PriceLevelsByPrices[order.Type][order.Price]
	if priceLevel == nil {
		priceLevel = ob.addNewPriceLevel(order)
	}

	priceLevel.AmountSum += order.Amount
	priceLevel.OrdersCount++

	plOrder := ob.OrdersByOrderID[order.ID] // this ID is an external ID
	if plOrder != nil {
		return nil
	}
	plOrder = &PriceLevelOrder{Order: order}
	ob.OrdersByOrderID[plOrder.Order.ID] = plOrder
	ob.OrdersCount[order.Type]++

	if priceLevel.OrderHead == nil {
		priceLevel.OrderHead = plOrder
		return plOrder
	}

	currPLOrder := priceLevel.OrderHead
	for {
		if order.BetterThan(currPLOrder.Order) {
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
	return plOrder
}

// AddOrder adds an order into the OrderBook.
func (ob *OrderBook) AddOrder(order Order) *TradeRequest {
	newPLOrder := ob.addNewPriceLevelOrder(order)
	if newPLOrder == nil {
		return nil
	}
	return ob.checksIfMatchs(newPLOrder)
}

// checksIfMatchs checks if order has a match.
func (ob *OrderBook) checksIfMatchs(plOrder *PriceLevelOrder) *TradeRequest {
	firstBuyPL := ob.PriceLevelsHeads[orders.OrderTypeBuy]
	firstSellPL := ob.PriceLevelsHeads[orders.OrderTypeSell]
	if firstBuyPL == nil || firstBuyPL.OrderHead == nil {
		return nil
	}
	if firstSellPL == nil || firstSellPL.OrderHead == nil {
		return nil
	}

	if firstBuyPL.OrderHead.Order.InTrade || firstSellPL.OrderHead.Order.InTrade {
		return nil
	}

	if firstBuyPL.OrderHead.Order.Price < firstSellPL.OrderHead.Order.Price {
		return nil
	}

	if !(firstBuyPL.OrderHead.Order.Mine || firstSellPL.OrderHead.Order.Mine) {
		return nil
	}

	// Match!
	amount := firstBuyPL.OrderHead.Order.Amount
	if firstBuyPL.OrderHead.Order.Amount > firstSellPL.OrderHead.Order.Amount {
		amount = firstSellPL.OrderHead.Order.Amount
	}

	var tradeReq *TradeRequest

	firstBuyPL.OrderHead.Order.InTrade = true
	firstSellPL.OrderHead.Order.InTrade = true

	if firstBuyPL.OrderHead.Order.Mine {
		tradeReq = &TradeRequest{
			InterestedOrder: firstBuyPL.OrderHead.Order,
			InterestOrder:   firstSellPL.OrderHead.Order,
			Amount:          amount,
		}
	} else {
		tradeReq = &TradeRequest{
			InterestedOrder: firstSellPL.OrderHead.Order,
			InterestOrder:   firstBuyPL.OrderHead.Order,
			Amount:          amount,
		}
	}
	if tradeReq != nil {
		log.Printf("Order match! %v-%v-$%v-%vqty: True, = %v-%v-$%v-%vqty (amount take %v)",
			tradeReq.InterestedOrder.Type, tradeReq.InterestedOrder.ID, tradeReq.InterestedOrder.Price, tradeReq.InterestedOrder.Amount,
			tradeReq.InterestOrder.Type, tradeReq.InterestOrder.ID, tradeReq.InterestOrder.Price, tradeReq.InterestOrder.Amount,
			tradeReq.Amount,
		)
	}
	return tradeReq
}

// DecOrderAmount decrement an order amount.
func (ob *OrderBook) DecOrderAmount(order Order) {
	plOrder := ob.OrdersByOrderID[order.ID]
	if plOrder == nil {
		return
	}
	plOrder.Order.InTrade = false
	plOrder.Order.Amount -= order.Amount
	if plOrder.Order.Amount <= 0 {
		ob.RemoveOrder(order)
	}
}

// RemoveOrder removes an order from the OrderBook.
func (ob *OrderBook) RemoveOrder(order Order) {
	priceLevel, _ := ob.PriceLevelsByPrices[order.Type][order.Price]
	if priceLevel == nil {
		return
	}

	plOrder := ob.OrdersByOrderID[order.ID]
	if plOrder == nil {
		return
	}

	if plOrder.Right != nil {
		plOrder.Right.Left = plOrder.Left
	}

	if plOrder.Left == nil { // head
		priceLevel.OrderHead = plOrder.Right
	} else {
		plOrder.Left.Right = plOrder.Right
	}

	plOrder.Left = nil
	plOrder.Right = nil

	delete(ob.OrdersByOrderID, order.ID)
	ob.OrdersCount[order.Type]--
	priceLevel.AmountSum -= order.Amount
	priceLevel.OrdersCount--
}

// GetBuyOrders returns a slice of buying orders.
func (ob *OrderBook) GetBuyOrders() []Order {
	orders := ob.getOrders(orders.OrderTypeBuy)
	return *orders
}

// GetSellOrders returns a slice of selling orders.
func (ob *OrderBook) GetSellOrders() []Order {
	orders := ob.getOrders(orders.OrderTypeSell)
	return *orders
}

// getOrders returns a slice of buying or selling orders.
func (ob *OrderBook) getOrders(orderType orders.OrderType) *[]Order {
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
