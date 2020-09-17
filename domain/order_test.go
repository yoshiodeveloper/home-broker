package domain_test

import (
	"home-broker/domain"
	"reflect"
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

func GetTestOrder(id domain.OrderID, externalOrderID string, amount int64, price domain.Money, Type string, exTimestamp time.Time) domain.Order {
	asset := GetTestAsset()
	return domain.Order{
		ID:                id,
		ExternalOrderID:   externalOrderID,
		UserID:            999,
		AssetID:           asset.ID,
		Amount:            amount,
		Price:             price,
		Type:              Type,
		ExternalTimestamp: exTimestamp,
		CreatedAt:         exTimestamp,
		UpdatedAt:         exTimestamp.Add(time.Hour * 2),
		DeletedAt:         time.Time{},
	}
}

func TestOrderAddOrder_AddBuyOrderOnEmptyPool_BuyOrderAdded(t *testing.T) {
	oType := domain.OrderTypeBuy
	order := GetTestOrder(999999, "ex999999", 99, domain.NewMoneyFromFloatString("99.9999999"), oType, baseTime)
	ob := domain.NewOrderBook(order.AssetID)
	if ob.BuyOrdersHead != nil {
		t.Errorf("BuyOrdersHead is not nil, expected a nil")
	}
	ob.AddOrder(&order)
	if ob.BuyOrdersHead == nil {
		t.Errorf("BuyOrdersHead is nil, expected as not nil")
	}
	if ob.BuyOrdersHead.Left != nil {
		t.Errorf("BuyOrdersHead.Left is not nil, expected as nil")
	}
	if ob.BuyOrdersHead.Right != nil {
		t.Errorf("BuyOrdersHead.Right is not nil, expected as nil")
	}
	if ob.BuyOrdersHead.Price != order.Price {
		t.Errorf("BuyOrdersHead.Price is %v, expected %v", ob.BuyOrdersHead.Price, order.Price)
	}
	if ob.BuyOrdersHead.AmountSum != order.Amount {
		t.Errorf("BuyOrdersHead.AmountSum is %v, expected %v", ob.BuyOrdersHead.AmountSum, order.Amount)
	}
	if ob.BuyOrdersHead.OrdersCount != 1 {
		t.Errorf("BuyOrdersHead.OrdersCount is %v, expected 1", ob.BuyOrdersHead.OrdersCount)
	}

	if ob.BuyOrdersHead.OrderHead == nil {
		t.Errorf("BuyOrdersHead.OrderHead is nil, expected not nil")
	}
	if ob.BuyOrdersHead.OrderHead.Left != nil {
		t.Errorf("BuyOrdersHead.OrderHead.Left is not nil, expected as nil")
	}
	if ob.BuyOrdersHead.OrderHead.Right != nil {
		t.Errorf("BuyOrdersHead.OrderHead.Right is not nil, expected as nil")
	}
	if !reflect.DeepEqual(ob.BuyOrdersHead.OrderHead.Order, order) {
		t.Errorf("ob.BuyOrdersHead.OrderHead.Order is different from order: %v != %v", ob.BuyOrdersHead.OrderHead.Order, order)
	}
}

func TestOrderAddOrder_AddSellOrderOnEmptyPool_SellOrderAdded(t *testing.T) {
	oType := domain.OrderTypeSell
	order := GetTestOrder(999999, "ex999999", 99, domain.NewMoneyFromFloatString("99.9999999"), oType, baseTime)
	ob := domain.NewOrderBook(order.AssetID)
	if ob.SellOrdersHead != nil {
		t.Errorf("SellOrdersHead is not nil, expected a nil")
	}
	ob.AddOrder(&order)
	if ob.SellOrdersHead == nil {
		t.Errorf("SellOrdersHead is nil, expected as not nil")
	}
	if ob.SellOrdersHead.Left != nil {
		t.Errorf("SellOrdersHead.Left is not nil, expected as nil")
	}
	if ob.SellOrdersHead.Right != nil {
		t.Errorf("SellOrdersHead.Right is not nil, expected as nil")
	}
	if ob.SellOrdersHead.Price != order.Price {
		t.Errorf("SellOrdersHead.Price is %v, expected %v", ob.SellOrdersHead.Price, order.Price)
	}
	if ob.SellOrdersHead.AmountSum != order.Amount {
		t.Errorf("SellOrdersHead.AmountSum is %v, expected %v", ob.SellOrdersHead.AmountSum, order.Amount)
	}
	if ob.SellOrdersHead.OrdersCount != 1 {
		t.Errorf("SellOrdersHead.OrdersCount is %v, expected 1", ob.SellOrdersHead.OrdersCount)
	}

	if ob.SellOrdersHead.OrderHead == nil {
		t.Errorf("SellOrdersHead.OrderHead is nil, expected not nil")
	}
	if ob.SellOrdersHead.OrderHead.Left != nil {
		t.Errorf("SellOrdersHead.OrderHead.Left is not nil, expected as nil")
	}
	if ob.SellOrdersHead.OrderHead.Right != nil {
		t.Errorf("SellOrdersHead.OrderHead.Right is not nil, expected as nil")
	}
	if !reflect.DeepEqual(ob.SellOrdersHead.OrderHead.Order, order) {
		t.Errorf("ob.SellOrdersHead.OrderHead.Order is different from order: %v != %v", ob.SellOrdersHead.OrderHead.Order, order)
	}
}

func TestOrderAddOrder_AddManyBuyOrders_BuyOrdersAddedSorted(t *testing.T) {
	oType := domain.OrderTypeBuy
	expectedOrder := []domain.Order{
		GetTestOrder(1, "ex1", 91, domain.NewMoneyFromFloatString("99.000005"), oType, baseTime), // < -- best bid
		GetTestOrder(2, "ex2", 97, domain.NewMoneyFromFloatString("99.000004"), oType, baseTime),
		GetTestOrder(3, "ex3", 93, domain.NewMoneyFromFloatString("99.000004"), oType, baseTime.Add(time.Millisecond)),
		GetTestOrder(4, "ex4", 95, domain.NewMoneyFromFloatString("99.000003"), oType, baseTime),
		GetTestOrder(5, "ex5", 95, domain.NewMoneyFromFloatString("99.000003"), oType, baseTime.Add(time.Millisecond)),
		GetTestOrder(6, "ex6", 98, domain.NewMoneyFromFloatString("99.000002"), oType, baseTime),
		GetTestOrder(7, "ex7", 96, domain.NewMoneyFromFloatString("99.000001"), oType, baseTime), // <-- worst bid
	}

	// Must be ordered by price and time.
	insertionOrder := []domain.Order{
		expectedOrder[5],
		expectedOrder[2],
		expectedOrder[6],
		expectedOrder[4],
		expectedOrder[0],
		expectedOrder[3],
		expectedOrder[1],
	}
	// GetTestOrder(6, "ex6", 97, domain.NewMoneyFromFloatString("99.000004"), oType, baseTime)
	// GetTestOrder(3, "ex3", 95, domain.NewMoneyFromFloatString("99.000003"), oType, baseTime.Add(time.Millisecond))
	// --
	// GetTestOrder(3, "ex3", 95, domain.NewMoneyFromFloatString("99.000003"), oType, baseTime.Add(time.Millisecond))
	// GetTestOrder(7, "ex7", 91, domain.NewMoneyFromFloatString("99.000005"), oType, baseTime)
	// GetTestOrder(6, "ex6", 97, domain.NewMoneyFromFloatString("99.000004"), oType, baseTime)

	ob := domain.NewOrderBook(expectedOrder[0].AssetID)

	for _, order := range insertionOrder {
		ob.AddOrder(&order)
	}

	if len(ob.GetBuyingOrders()) != len(expectedOrder) {
		t.Errorf("BuyOrders does not has %d itens: %v", len(expectedOrder), len(ob.GetBuyingOrders()))
	}
	for i, order := range expectedOrder {
		t.Logf("verifying item in index %v", i)
		if !reflect.DeepEqual(ob.GetBuyingOrders()[i], order) {
			t.Errorf("ob.BuyOrdersHead.OrderHead.Order[%v] is different from order: %v != %v", i, ob.GetBuyingOrders()[i], order)
		}
	}
}
