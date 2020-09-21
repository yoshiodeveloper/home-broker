package orderbooksgin

import (
	"fmt"
	"home-broker/assets"
	"home-broker/core"
	"home-broker/orderbooks"
	"home-broker/orders"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	apiErrorInvalidJSON = core.NewAPIError("Invalid JSON.", 400)
)

// OrderBookController represents an order controller.
type OrderBookController struct {
	assetID assets.AssetID
	uc      orderbooks.OrderBookUseCases
}

// NewOrderBookController creates a new OrderBookController.
func NewOrderBookController(assetID assets.AssetID, uc orderbooks.OrderBookUseCases) OrderBookController {
	return OrderBookController{assetID: assetID, uc: uc}
}

// Webhook receives the updates in the order book.
func (orderBookC OrderBookController) Webhook(c *gin.Context) {
	assetID := assets.AssetID(c.Param("asset_id"))
	if orderBookC.assetID != assetID {
		c.Error(core.NewAPIError(fmt.Sprintf("Check the URL. This host only handles orders of asset \"%v\".", orderBookC.assetID), 400))
		return
	}
	var json orders.ExternalUpdate
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(apiErrorInvalidJSON)
		return
	}
	response, err := orderBookC.uc.Webhook(json)
	if err != nil {
		errVal, ok := err.(core.ErrValidation)
		if ok {
			c.Error(core.NewAPIErrorFromErrValidation(errVal))
			return
		}
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response)
}
