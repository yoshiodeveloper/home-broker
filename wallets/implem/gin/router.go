package walletsgin

import (
	"home-broker/wallets"

	"github.com/gin-gonic/gin"
)

// WalletRouter represents a wallets router.
type WalletRouter struct {
	uc wallets.WalletUseCases
}

// NewWalletRouter creates a new Router.
func NewWalletRouter(uc wallets.WalletUseCases) WalletRouter {
	return WalletRouter{uc: uc}
}

// SetupRouter setups wallets router.
func (wr WalletRouter) SetupRouter(router *gin.Engine) {
	walletC := NewWalletController(wr.uc)
	router.GET("/api/v1/wallet/:user_id/", walletC.GetWallet)
}
