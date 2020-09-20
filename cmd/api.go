package cmd

import (
	"home-broker/core/implem/postgresql"
	"log"

	coregin "home-broker/core/implem/gin"

	"home-broker/users"
	userspostgresql "home-broker/users/implem/postgresql"
	"home-broker/wallets"
	walletsginserver "home-broker/wallets/implem/gin"
	walletspostgresql "home-broker/wallets/implem/postgresql"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API webserver",
	Run:   runAPIServer,
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runAPIServer(cmd *cobra.Command, args []string) {
	mainDB := postgresql.NewDB("localhost", 5432, "homebroker", "123456", "homebroker")
	err := mainDB.Open()
	if err != nil {
		log.Panic(err)
	}

	mainDB.GetDB().AutoMigrate()

	userDB := userspostgresql.NewUserDB(mainDB)
	walletDB := walletspostgresql.NewWalletDB(mainDB)

	userUC := users.NewUserUseCases(userDB)
	walletUC := wallets.NewWalletUseCases(walletDB, userUC)

	router := gin.Default()
	router.Use(coregin.MiddlewareAPIError())

	walletrouter := walletsginserver.NewWalletRouter(walletUC)
	walletrouter.SetupRouter(router)

	router.Run()
}
