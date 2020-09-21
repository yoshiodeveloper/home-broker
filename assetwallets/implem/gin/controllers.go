package walletsgin

import (
	"home-broker/assets"
	"home-broker/assetwallets"
	"home-broker/core"
	"home-broker/users"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	apiErrorInvalidJSON   = core.NewAPIError("Invalid JSON.", 400)
	apiErrorInvalidUserID = core.NewAPIError("Invalid user ID.", 400)
)

// AssetWalletController represents a wallet controller.
type AssetWalletController struct {
	uc assetwallets.AssetWalletUseCases
}

// NewAssetWalletController creates a new WalletController.
func NewAssetWalletController(uc assetwallets.AssetWalletUseCases) AssetWalletController {
	return AssetWalletController{uc: uc}
}

// GetAssetWallet returns an user asset wallet.
func (assetWalletC AssetWalletController) GetAssetWallet(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.Error(apiErrorInvalidUserID)
		return
	}
	assetID := c.Param("asset_id")
	entity, _, _, err := assetWalletC.uc.GetAssetWallet(users.UserID(userID), assets.AssetID(assetID))
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
