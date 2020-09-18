package domain_test

import (
	"fmt"
	"home-broker/domain"
	"math/rand"
	"testing"
	"time"
)

var (
	baseTime = time.Date(2020, time.Month(1), 10, 11, 12, 13, 14, time.UTC)
)

func GetTestAsset() domain.Asset {
	return domain.Asset{
		ID:        "PETR4",
		Name:      "PETROBRAS",
		CreatedAt: baseTime,
		UpdatedAt: baseTime.Add(time.Hour * 2),
		DeletedAt: time.Time{},
	}
}

func GetTestOrder(id int64, orderType domain.OrderType, price domain.Money, amount int64, exTimestamp time.Time) domain.Order {
	asset := GetTestAsset()
	var o domain.Order
	if orderType == domain.OrderTypeBuy {
		o = domain.NewBuyOrder(asset.ID, amount, price)
	} else {
		o = domain.NewSellOrder(asset.ID, amount, price)
	}
	o.ID = domain.OrderID(id)
	o.ExternalID = domain.ExternalOrderID(fmt.Sprintf("ex%d", id))
	o.ExternalTimestamp = exTimestamp
	o.CreatedAt = exTimestamp
	o.UpdatedAt = exTimestamp.Add(time.Hour * 2)
	o.DeletedAt = time.Time{}
	o.UserID = 999
	return o
}

func TestOrderAddOrder_AddManyBuyOrders_BuyOrdersAddedSorted(t *testing.T) {
	tests := 1000
	expectedOrders := make([]domain.Order, 0)
	insertionOrders := make([]domain.Order, 0)
	exTime := baseTime

	// This creates a list of offers in the expected order.
	// Each price level will have 10 orders with the same price but different time.
	price := domain.Money(tests)
	for i := tests - 1; i >= 0; i-- {
		if i%5 == 0 {
			// Creates 10 equals prices (same price level).
			price -= domain.Money(1) // decrementing
		}
		order := GetTestOrder(int64(i+1), domain.OrderTypeBuy, price, int64(i+1), exTime)
		expectedOrders = append(expectedOrders, order)
		insertionOrders = append(insertionOrders, order)
		if i%2 == 0 {
			// The time is the same every even.
			exTime = exTime.Add(time.Nanosecond)
		}
	}

	t.Logf("%d order(s) was created for the test", len(expectedOrders))

	rand.Seed(1) // Shuffes always in the same order.
	rand.Shuffle(len(insertionOrders), func(i, j int) {
		insertionOrders[i], insertionOrders[j] = insertionOrders[j], insertionOrders[i]
	})

	ob := domain.NewOrderBook(expectedOrders[0].AssetID)
	for i, order := range insertionOrders {
		if order.Price == 0 {
			t.Errorf("insertionOrders[%d] has no Price: %v", i, order)
		}
		if order.Type == "" {
			t.Errorf("insertionOrders[%d] has no Type: %v", i, order)
		}
		ob.AddOrder(order)
	}

	orders := ob.GetBuyOrders()
	if len(orders) != len(expectedOrders) {
		t.Errorf("orders has %d itens, expected %d", len(orders), len(expectedOrders))
	}

	for i, expectedOrder := range expectedOrders {
		order := orders[i]
		// We shouldn't compare IDs because when the price and timestamp are the
		// same the items don't have a specific order.
		if order.Price != expectedOrder.Price {
			t.Errorf("order[%d].Price is %v, expected %v", i, order.Price, expectedOrder.Price)
		}
		if order.ExternalTimestamp != expectedOrder.ExternalTimestamp {
			t.Errorf("order[%d].ExternalTimestamp is %v, expected %v", i, order.ExternalTimestamp, expectedOrder.ExternalTimestamp)
		}
	}
}

func TestOrderAddOrder_AddManySellOrders_SellOrdersAddedSorted(t *testing.T) {
	tests := 1000
	expectedOrders := make([]domain.Order, 0)
	insertionOrders := make([]domain.Order, 0)
	exTime := baseTime

	// This creates a list of offers in the expected order.
	// Each price level will have 10 orders with the same price but different time.
	price := domain.Money(0)
	for i := 0; i < tests; i++ {
		if i%5 == 0 {
			// Creates 10 equals prices (same price level).
			price += domain.Money(1) // incrementing
		}
		order := GetTestOrder(int64(i+1), domain.OrderTypeSell, price, int64(i+1), exTime)
		expectedOrders = append(expectedOrders, order)
		insertionOrders = append(insertionOrders, order)
		if i%2 == 0 {
			// The time is the same every even.
			exTime = exTime.Add(time.Nanosecond)
		}
	}

	t.Logf("%d order(s) was created for the test", len(expectedOrders))

	rand.Seed(1) // Shuffes always in the same order.
	rand.Shuffle(len(insertionOrders), func(i, j int) {
		insertionOrders[i], insertionOrders[j] = insertionOrders[j], insertionOrders[i]
	})

	ob := domain.NewOrderBook(expectedOrders[0].AssetID)
	for i, order := range insertionOrders {
		if order.Price == 0 {
			t.Errorf("insertionOrders[%d] has no Price: %v", i, order)
		}
		if order.Type == "" {
			t.Errorf("insertionOrders[%d] has no Type: %v", i, order)
		}
		ob.AddOrder(order)
	}

	orders := ob.GetSellOrders()
	if len(orders) != len(expectedOrders) {
		t.Errorf("orders has %d itens, expected %d", len(orders), len(expectedOrders))
	}

	for i, expectedOrder := range expectedOrders {
		order := orders[i]
		// We shouldn't compare IDs because when the price and timestamp are the
		// same the items don't have a specific order.
		if order.Price != expectedOrder.Price {
			t.Errorf("order[%d].Price is %v, expected %v", i, order.Price, expectedOrder.Price)
		}
		if order.ExternalTimestamp != expectedOrder.ExternalTimestamp {
			t.Errorf("order[%d].ExternalTimestamp is %v, expected %v", i, order.ExternalTimestamp, expectedOrder.ExternalTimestamp)
		}
	}
}

func BenchmarkOderBookInsertion(b *testing.B) {
	asset := GetTestAsset()
	ob := domain.NewOrderBook(asset.ID)
	var orderType domain.OrderType
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			orderType = domain.OrderTypeBuy
		} else {
			orderType = domain.OrderTypeSell
		}
		orderID := rand.Int63() + 1
		amount := rand.Int63n(10000) + 1
		exTime := time.Now()
		// 10k price levels (ex: $0.01 ~ $100.01)
		price := domain.Money(rand.Int63n(10000) + 1)
		order := GetTestOrder(orderID, orderType, price, amount, exTime)
		ob.AddOrder(order)
	}
}
