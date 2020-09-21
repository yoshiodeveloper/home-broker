package ordersgin

import (
	"home-broker/assets"
	"home-broker/core"
	"home-broker/money"
	"home-broker/orders"
	"home-broker/users"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	apiErrorInvalidJSON    = core.NewAPIError("Invalid JSON.", 400)
	apiErrorNoFunds        = core.NewAPIError("No funds available.", 400)
	apiErrorInvalidPrice   = core.NewAPIError("Invalid price.", 400)
	apiErrorInvalidAmount  = core.NewAPIError("Invalid amount.", 400)
	apiErrorInvalidOrderID = core.NewAPIError("Invalid order ID.", 400)
	apiErrorInvalidUserID  = core.NewAPIError("Invalid user ID.", 400)
	apiErrorInvalidAssetID = core.NewAPIError("Invalid asset ID.", 400)
)

// OrderController represents an order controller.
type OrderController struct {
	uc orders.OrderUseCases
}

// NewOrderController creates a new WalletController.
func NewOrderController(uc orders.OrderUseCases) OrderController {
	return OrderController{uc: uc}
}

// AddOrderJSON is the JSON received on AddBuyOrder or AddSellOrder.
type AddOrderJSON struct {
	UserID  users.UserID     `json:"user_id"`
	AssetID assets.AssetID   `json:"asset_id"`
	Price   money.Money      `json:"price"`
	Amount  assets.AssetUnit `json:"amount"`
}

// GetOrder returns an order.
func (orderC OrderController) GetOrder(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("order_id"), 10, 64)
	if err != nil {
		c.Error(apiErrorInvalidOrderID)
		return
	}
	entity, err := orderC.uc.GetOrder(orders.OrderID(orderID))
	if err != nil {
		c.Error(err)
		return
	}
	if entity == nil {
		c.Error(core.NewAPIError("Not found", 404))
		return
	}
	c.JSON(http.StatusOK, entity)
}

// BuyOrder adds an buy Order.
func (orderC OrderController) BuyOrder(c *gin.Context) {
	var json AddOrderJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(apiErrorInvalidJSON)
		return
	}
	entity, err := orderC.uc.BuyOrder(json.UserID, json.AssetID, json.Price, json.Amount)
	if err != nil {
		errVal, ok := err.(core.ErrValidation)
		if ok {
			c.Error(core.NewAPIErrorFromErrValidation(errVal))
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, entity)
}

// SellOrder adds an buy Order.
func (orderC OrderController) SellOrder(c *gin.Context) {
	var json AddOrderJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(apiErrorInvalidJSON)
		return
	}
	entity, err := orderC.uc.SellOrder(json.UserID, json.AssetID, json.Price, json.Amount)
	if err != nil {
		errVal, ok := err.(core.ErrValidation)
		if ok {
			c.Error(core.NewAPIErrorFromErrValidation(errVal))
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, entity)
}

// CancelOrder returns an order.
func (orderC OrderController) CancelOrder(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("order_id"), 10, 64)
	if err != nil {
		c.Error(apiErrorInvalidOrderID)
		return
	}
	entity, err := orderC.uc.CancelOrder(orders.OrderID(orderID))
	if err != nil {
		c.Error(err)
		return
	}
	if entity == nil {
		c.Error(core.NewAPIError("Not found", 404))
		return
	}
	c.JSON(http.StatusOK, entity)
}
