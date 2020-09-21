package cmd

import (
	"fmt"
	"home-broker/assets"
	"home-broker/config"
	coregin "home-broker/core/implem/gin"
	"home-broker/orderbooks"
	orderbooksgin "home-broker/orderbooks/implem/gin"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// orderbookCmd represents the orderbook command
var orderbookCmd = &cobra.Command{
	Use:   "orderbook",
	Short: "Starts the order book service",
	Long:  `Important: You must have only one instance of the order book service.`,
	Run:   startOrderBook,
}

func init() {
	rootCmd.AddCommand(orderbookCmd)
	orderbookCmd.Flags().String("asset", "VIBR", "The asset ID this Order Book must handle.")
}

func startOrderBook(cmd *cobra.Command, args []string) {
	ginConfig := config.NewGinConfigFromViper(viper.GetViper())
	assetID, err := cmd.Flags().GetString("asset")
	if err != nil {
		log.Fatal(err)
	}

	orderBook := orderbooks.NewOrderBook(assets.AssetID(assetID))
	orderBookUC := orderbooks.NewOrderBookUseCases(orderBook)

	router := gin.Default()
	router.Use(coregin.MiddlewareAPIError())

	orderBookRouter := orderbooksgin.NewOrderBookRouter(orderBook.AssetID, orderBookUC)
	orderBookRouter.SetupRouter(router)

	log.Printf("\n\n#\n# IMPORTANT: You must execute only one instance of the Order Book for asset \"%v\"\n#\n", assetID)
	router.Run(fmt.Sprintf(":%d", ginConfig.Port))
}
