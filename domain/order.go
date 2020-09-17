package domain

import (
	"fmt"
	"time"
)

const (
	// OrderTypeBuy is an order type for a buy.
	OrderTypeBuy = "buy"

	// OrderTypeSell is an Order type for a sell.
	OrderTypeSell = "sell"
)

// OrderID represents the Order ID type.
//   This eases a future DB change.
type OrderID int64

// Order is an entity for buying or selling intentions.
type Order struct {
	ID                OrderID
	ExternalOrderID   string // ID from external service (as Stock Exchange)
	UserID            UserID
	AssetID           AssetID
	Amount            int64
	Price             Money
	Type              string
	ExternalTimestamp time.Time // Date/Time from external service.
	CreatedAt         time.Time // Date/Time in this service.
	UpdatedAt         time.Time
	DeletedAt         time.Time
}

// BetterThan returns true if this order is better than order parameter.
// Both orders must be the same Type.
func (o *Order) BetterThan(order *Order) bool {
	if o.Type == OrderTypeBuy && o.Price > order.Price {
		return true
	} else if o.Type == OrderTypeSell && o.Price < order.Price {
		return true
	}
	if o.Price == order.Price && o.CreatedAt.Before(order.CreatedAt) {
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
	Left  *PriceLevel
	Right *PriceLevel
	//Next      *PriceLevel
	OrderHead *PriceLevelOrder
	Price     Money
	Type      string
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

	// OrdersByOrderID maps an Order ID to a PriceLevelOrder.
	// This improves the search of a order inside the order book at O(1).
	// The downside is the rehash of a hashmap.
	// This can be improved with a AVL BTree.
	OrdersByOrderID map[OrderID]*PriceLevelOrder

	// PriceLevelsBuyByPrice maps a Price to a OrderBookPriceLevels.
	// This improves the search of a order inside the order book at O(1).
	// The downside is the rehash of a hashmap.
	// This can be improved with a AVL BTree.
	PriceLevelsBuyByPrice map[Money]*PriceLevel

	// PriceLevelsSellByPrice maps a Price to a OrderBookPriceLevels.
	PriceLevelsSellByPrice map[Money]*PriceLevel

	// BuyOrdersHead is a linked list of buying orders.
	// The best buying orders will be the first node (head).
	BuyOrdersHead *PriceLevel

	BuyOrdersCount int64

	// SellOrdersHead is a linked list of buying orders.
	// The best selling orders will be the first node (head).
	SellOrdersHead *PriceLevel

	SellOrdersCount int64
}

// NewOrderBook creates a new OrderBook.
func NewOrderBook(assetID AssetID) *OrderBook {
	return &OrderBook{
		AssetID:                assetID,
		OrdersByOrderID:        make(map[OrderID]*PriceLevelOrder),
		PriceLevelsBuyByPrice:  make(map[Money]*PriceLevel),
		PriceLevelsSellByPrice: make(map[Money]*PriceLevel),
	}
}

// AddPriceLevel adds a PriceLevel into the OrderBook.
// Nothing will be done if it already exists.
func (ob *OrderBook) AddPriceLevel(orderType string, price Money) *PriceLevel {
	var plsByPrice *map[Money]*PriceLevel
	var head **PriceLevel

	switch orderType {
	case OrderTypeBuy:
		plsByPrice = &ob.PriceLevelsBuyByPrice
		head = &ob.BuyOrdersHead
	case OrderTypeSell:
		plsByPrice = &ob.PriceLevelsSellByPrice
		head = &ob.SellOrdersHead
	default:
		panic(fmt.Errorf("programming error. %v is not a mapped type", orderType))
	}

	priceLevel, hasKey := (*plsByPrice)[price]
	if hasKey {
		return priceLevel
	}

	priceLevel = &PriceLevel{Price: price, Type: orderType}

	if *head == nil {
		*head = priceLevel
	} else {
		currPL := *head
		for {
			if priceLevel.BetterOfferThan(currPL) {
				if currPL.Left == nil {
					// currPL is the first node.
					*head = priceLevel
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

	}

	(*plsByPrice)[priceLevel.Price] = priceLevel
	return priceLevel
}

// AddOrder adds a order into the OrderBook.
func (ob *OrderBook) AddOrder(order *Order) {
	// Order was copied.
	// We are not using pointer as an Order can not be changed by reference inside an OrderBook.
	plOrder := &PriceLevelOrder{Order: *order}

	priceLevel := ob.AddPriceLevel(order.Type, order.Price)
	priceLevel.AmountSum += order.Amount
	priceLevel.OrdersCount++

	if priceLevel.OrderHead == nil {
		priceLevel.OrderHead = plOrder
	} else {
		currPLOrder := priceLevel.OrderHead
		for {
			if order.BetterThan(&currPLOrder.Order) {
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

	ob.OrdersByOrderID[plOrder.Order.ID] = plOrder

	switch plOrder.Order.Type {
	case OrderTypeBuy:
		ob.BuyOrdersCount++
	case OrderTypeSell:
		ob.SellOrdersCount++
	default:
		panic(fmt.Errorf("programming error. %v is not a mapped type", plOrder.Order.Type))
	}
}

// GetBuyingOrders returns a slice of buying orders.
func (ob *OrderBook) GetBuyingOrders() []Order {
	orders := make([]Order, 0, ob.BuyOrdersCount)
	currPriceLevel := ob.BuyOrdersHead
	for currPriceLevel != nil {
		currPLOrder := currPriceLevel.OrderHead
		for currPLOrder != nil {
			orders = append(orders, currPLOrder.Order)
			currPLOrder = currPLOrder.Right
		}
		currPriceLevel = currPriceLevel.Right
	}
	return orders
}

// RemoveOrder removes an order into the order book.
// func (ob *OrderBook) RemoveOrder(orderID OrderID)
