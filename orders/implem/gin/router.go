package ordersgin

import (
	"home-broker/orders"

	"github.com/gin-gonic/gin"
)

// OrderRouter represents an orders router.
type OrderRouter struct {
	uc orders.OrderUseCases
}

// NewOrderRouter creates a new Router.
func NewOrderRouter(uc orders.OrderUseCases) OrderRouter {
	return OrderRouter{uc: uc}
}

// SetupRouter setups orders router.
func (wr OrderRouter) SetupRouter(router *gin.Engine) {
	orderC := NewOrderController(wr.uc)
	v1 := router.Group("/api/v1/orders")
	{
		v1.POST("buy/", orderC.BuyOrder)
		v1.POST("sell/", orderC.SellOrder)
		v1.GET(":order_id/", orderC.GetOrder)
		v1.DELETE(":order_id/", orderC.CancelOrder)
	}
}
