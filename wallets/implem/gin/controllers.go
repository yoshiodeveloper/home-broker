package walletsgin

import (
	"errors"
	"home-broker/core"
	"home-broker/money"
	"home-broker/users"
	"home-broker/wallets"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	apiErrorInvalidJSON        = core.NewAPIError("Invalid JSON.", 400)
	apiErrorInvalidFundsAmount = core.NewAPIError("The funds amount is invalid.", 400)
	apiErrorInvalidUserID      = core.NewAPIError("User ID is invalid.", 400)
)

// WalletController represents a wallet controller.
type WalletController struct {
	uc wallets.WalletUseCases
}

// NewWalletController creates a new WalletController.
func NewWalletController(uc wallets.WalletUseCases) WalletController {
	return WalletController{uc: uc}
}

// GetWallet returns an user wallet.
func (walletC WalletController) GetWallet(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.Error(apiErrorInvalidUserID)
		return
	}
	//entity, created, userCreated, err = walletUC.GetWallet(domain.UserID(userID))
	entity, _, _, err := walletC.uc.GetWallet(users.UserID(userID))
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, entity)
}

// AddFundsJSON is the JSON received on AddFunds.
type AddFundsJSON struct {
	Amount money.Money
}

// AddFunds adds funds to an user wallet.
func (walletC WalletController) AddFunds(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.Error(apiErrorInvalidUserID)
		return
	}

	var json AddFundsJSON
	if err := c.ShouldBindJSON(&json); err != nil {
		c.Error(apiErrorInvalidJSON)
		return
	}

	entity, err := walletC.uc.AddFunds(users.UserID(userID), json.Amount)
	if err != nil {
		if errors.Is(err, wallets.ErrInvalidFundsAmount) {
			c.Error(apiErrorInvalidFundsAmount)
			return
		}
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, entity)
}
