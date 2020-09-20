package cmd

import (
	"home-broker/config"
	"home-broker/core/implem/postgresql"
	userspostgresql "home-broker/users/implem/postgresql"
	walletspostgresql "home-broker/wallets/implem/postgresql"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply migration on the database.",
	Run: func(cmd *cobra.Command, args []string) {
		pgConfig := config.NewPostgreSQLConfigFromViper(viper.GetViper())
		mainDB := postgresql.NewDB(pgConfig.Host, pgConfig.Port, pgConfig.User, pgConfig.Password, pgConfig.Name)
		err := mainDB.Open()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("applying UserModel...")
		mainDB.GetDB().AutoMigrate(&userspostgresql.UserModel{})
		log.Println("applying WalletModel...")
		mainDB.GetDB().AutoMigrate(&walletspostgresql.WalletModel{})
		log.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
