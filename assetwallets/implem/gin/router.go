package walletsgin

import (
	"home-broker/assetwallets"

	"github.com/gin-gonic/gin"
)

// AssetWalletRouter represents a wallets router.
type AssetWalletRouter struct {
	uc assetwallets.AssetWalletUseCases
}

// NewAssetWalletRouter creates a new Router.
func NewAssetWalletRouter(uc assetwallets.AssetWalletUseCases) AssetWalletRouter {
	return AssetWalletRouter{uc: uc}
}

// SetupRouter setups asset wallets router.
func (wr AssetWalletRouter) SetupRouter(router *gin.Engine) {
	assetWalletC := NewAssetWalletController(wr.uc)
	v1 := router.Group("/api/v1/wallets")
	{
		v1.GET(":user_id/assets/:asset_id/", assetWalletC.GetAssetWallet)
	}
}
