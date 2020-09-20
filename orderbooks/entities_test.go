package orderbooks_test

import (
	"home-broker/money"
	"home-broker/orderbooks"
	"home-broker/orders"
	testsassets "home-broker/tests/assets"
	testsorders "home-broker/tests/orders"
	"math/rand"
	"testing"
	"time"
)

func TestOrderAddOrder_AddManyBuyOrders_BuyOrdersAddedSorted(t *testing.T) {
	tests := 1000
	expectedOrders := make([]orders.Order, 0)
	insertionOrders := make([]orders.Order, 0)
	exTime := testsorders.BaseTime

	// This creates a list of offers in the expected order.
	// Each price level will have 10 orders with the same price but different time.
	price := money.Money(tests)
	for i := tests - 1; i >= 0; i-- {
		if i%5 == 0 {
			// Creates 10 equals prices (same price level).
			price -= money.Money(1) // decrementing
		}
		order := testsorders.GetOrder(orders.OrderID(i+1), orders.OrderTypeBuy, price, int64(i+1), exTime)
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

	ob := orderbooks.NewOrderBook(expectedOrders[0].AssetID)
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
	expectedOrders := make([]orders.Order, 0)
	insertionOrders := make([]orders.Order, 0)
	exTime := testsorders.BaseTime

	// This creates a list of offers in the expected order.
	// Each price level will have 10 orders with the same price but different time.
	price := money.Money(0)
	for i := 0; i < tests; i++ {
		if i%5 == 0 {
			// Creates 10 equals prices (same price level).
			price += money.Money(1) // incrementing
		}
		order := testsorders.GetOrder(orders.OrderID(i+1), orders.OrderTypeSell, price, int64(i+1), exTime)
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

	ob := orderbooks.NewOrderBook(expectedOrders[0].AssetID)
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
	asset := testsassets.GetAsset()
	ob := orderbooks.NewOrderBook(asset.ID)
	var orderType orders.OrderType
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			orderType = orders.OrderTypeBuy
		} else {
			orderType = orders.OrderTypeSell
		}
		orderID := orders.OrderID(rand.Int63() + 1)
		amount := rand.Int63n(10000) + 1
		exTime := time.Now()
		// 100k price levels (ex: $0.01 ~ $1000.01)
		price := money.Money(rand.Int63n(100000) + 1)
		order := testsorders.GetOrder(orderID, orderType, price, amount, exTime)
		ob.AddOrder(order)
	}
}

func BenchmarkOderBookSnapshots(b *testing.B) {
	asset := testsassets.GetAsset()
	ob := orderbooks.NewOrderBook(asset.ID)
	var orderType orders.OrderType
	for i := 0; i < 100000; i++ {
		if i%2 == 0 {
			orderType = orders.OrderTypeBuy
		} else {
			orderType = orders.OrderTypeSell
		}
		orderID := orders.OrderID(rand.Int63() + 1)
		amount := rand.Int63n(10000) + 1
		exTime := time.Now()
		// 100k price levels (ex: $0.01 ~ $1000.01)
		price := money.Money(rand.Int63n(100000) + 1)
		order := testsorders.GetOrder(orderID, orderType, price, amount, exTime)
		ob.AddOrder(order)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ob.GetBuyOrders()
		_ = ob.GetSellOrders()
	}
}
