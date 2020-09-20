package tests

import (
	"fmt"
	"home-broker/money"
	"home-broker/orders"
	testsassets "home-broker/tests/assets"
	"time"
)

var (
	// BaseTime is a base time for orders.
	BaseTime = time.Date(2020, time.Month(1), 10, 11, 12, 13, 14, time.UTC)
)

// GetOrder returns an order entity.
func GetOrder(id orders.OrderID, orderType orders.OrderType, price money.Money, amount int64, exTimestamp time.Time) orders.Order {
	asset := testsassets.GetAsset()
	var o orders.Order
	if orderType == orders.OrderTypeBuy {
		o = orders.NewBuyOrder(asset.ID, amount, price)
	} else {
		o = orders.NewSellOrder(asset.ID, amount, price)
	}
	o.ID = orders.OrderID(id)
	o.ExternalID = orders.ExternalOrderID(fmt.Sprintf("ex%d", id))
	o.ExternalTimestamp = exTimestamp
	o.CreatedAt = exTimestamp
	o.UpdatedAt = exTimestamp.Add(time.Hour * 2)
	o.DeletedAt = time.Time{}
	o.UserID = 999
	return o
}
