package coregin

/*
// Router represents a router.
type Router struct {
	useCases usecases.UseCases
}

// NewRouter creates a new Router.
func NewRouter(useCases usecases.UseCases) Router {
	return Router{useCases: useCases}
}

func (r Router) getUseCases() usecases.UseCases {
	return r.useCases
}

// GetGinRouter return a gin Router.
func (r Router) GetGinRouter() *gin.Engine {
	// dbClient := postgresql.NewDBClient("localhost", 5432, "homebroker", "123456", "homebroker")
	// useCases := usecases.NewUseCases(r.getRepositories())

	controllers := NewControllers(r.useCases)
	walletController := controllers.GetWalletController()

	router := gin.Default()

	router.Use(MiddlewareAPIError())

	router.GET("/api/v1/wallet/:user_id/", walletController.GetWallet)

	return router
}
*/
