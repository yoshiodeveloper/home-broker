package walletsgin

import (
	"home-broker/core"
	"home-broker/users"
	"home-broker/wallets"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		apierror := core.NewAPIError("User ID is not valid.", 400)
		c.Error(apierror)
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
