package orderbooksgin

import (
	"home-broker/assets"
	"home-broker/orderbooks"

	"github.com/gin-gonic/gin"
)

// OrderBookRouter represents an orders router.
type OrderBookRouter struct {
	assetID assets.AssetID
	uc      orderbooks.OrderBookUseCases
}

// NewOrderBookRouter creates a new Router.
func NewOrderBookRouter(assetID assets.AssetID, uc orderbooks.OrderBookUseCases) OrderBookRouter {
	return OrderBookRouter{assetID: assetID, uc: uc}
}

// SetupRouter setups orders router.
func (wr OrderBookRouter) SetupRouter(router *gin.Engine) {

	orderBookC := NewOrderBookController(wr.assetID, wr.uc)
	v1 := router.Group("/api/v1/orderbooks")
	{
		v1.POST(":asset_id/webhook/", orderBookC.Webhook)
	}
}
