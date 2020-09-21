package cmd

import (
	"fmt"
	"home-broker/assetwallets"
	"home-broker/config"
	"home-broker/core/implem/postgresql"
	"home-broker/orders"
	"log"
	"strings"

	"github.com/spf13/viper"

	coregin "home-broker/core/implem/gin"

	"home-broker/users"
	userspostgresql "home-broker/users/implem/postgresql"
	"home-broker/wallets"

	walletsginserver "home-broker/wallets/implem/gin"
	walletspostgresql "home-broker/wallets/implem/postgresql"

	assetwalletsginserver "home-broker/assetwallets/implem/gin"
	assetwalletspostgresql "home-broker/assetwallets/implem/postgresql"

	ordersginserver "home-broker/orders/implem/gin"
	orderspostgresql "home-broker/orders/implem/postgresql"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd *cobra.Command

func init() {
	apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Start the API webserver",
		Run:   runAPIServer,
	}
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	apiCmd.Flags().StringP("orderbook-host", "o", "http://localhost:8081", "The order book host API for a specific asset. (ex: \"http://localhost:8081\"")
	apiCmd.MarkFlagRequired("orderbook-host")
}

func runAPIServer(cmd *cobra.Command, args []string) {
	pgConfig := config.NewPostgreSQLConfigFromViper(viper.GetViper())
	ginConfig := config.NewGinConfigFromViper(viper.GetViper())

	orderBookHost, err := apiCmd.Flags().GetString("orderbook-host")
	if err != nil {
		log.Fatal(err)
	}
	if !strings.HasPrefix(orderBookHost, "http") {
		log.Fatal("The order book host must have http:// or https://")
	}

	mainDB := postgresql.NewDB(pgConfig.Host, pgConfig.Port, pgConfig.User, pgConfig.Password, pgConfig.Name)
	err = mainDB.Open()
	if err != nil {
		log.Fatal(err)
	}

	mainDB.GetDB().AutoMigrate()

	userDB := userspostgresql.NewUserDB(mainDB)
	walletDB := walletspostgresql.NewWalletDB(mainDB)
	assetWalletDB := assetwalletspostgresql.NewAssetWalletDB(mainDB)
	orderDB := orderspostgresql.NewOrderDB(mainDB)

	userUC := users.NewUserUseCases(userDB)
	walletUC := wallets.NewWalletUseCases(walletDB, userUC)
	assetWalletUC := assetwallets.NewAssetWalletUseCases(assetWalletDB, userUC)
	orderUC := orders.NewOrderUseCases(orderDB, walletUC, assetWalletUC, orderBookHost)

	if ginConfig.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(coregin.MiddlewareAPIError())

	walletRouter := walletsginserver.NewWalletRouter(walletUC)
	walletRouter.SetupRouter(router)

	assetwalletRouter := assetwalletsginserver.NewAssetWalletRouter(assetWalletUC)
	assetwalletRouter.SetupRouter(router)

	orderRouter := ordersginserver.NewOrderRouter(orderUC)
	orderRouter.SetupRouter(router)

	router.Run(fmt.Sprintf(":%d", ginConfig.Port))
}
