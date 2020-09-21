package cmd

import (
	"home-broker/assets"
	assetspostgresql "home-broker/assets/implem/postgresql"
	"home-broker/config"
	"home-broker/core/implem/postgresql"

	assetwalletspostgresql "home-broker/assetwallets/implem/postgresql"
	orderspostgresql "home-broker/orders/implem/postgresql"
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

		log.Println("applying AssetModel...")
		mainDB.GetDB().AutoMigrate(&assetspostgresql.AssetModel{})

		log.Println("applying UserModel...")
		mainDB.GetDB().AutoMigrate(&userspostgresql.UserModel{})

		log.Println("applying WalletModel...")
		mainDB.GetDB().AutoMigrate(&walletspostgresql.WalletModel{})

		log.Println("applying AssetWalletModel...")
		mainDB.GetDB().AutoMigrate(&assetwalletspostgresql.AssetWalletModel{})

		log.Println("applying OrderModel...")
		mainDB.GetDB().AutoMigrate(&orderspostgresql.OrderModel{})

		log.Println("inserting initial Assets data...")
		assets := []assets.Asset{
			assets.Asset{ID: "VIBR", Name: "Vibranium", ExchangeID: "VIBR"},
		}
		assetDB := assetspostgresql.NewAssetDB(mainDB)
		for _, asset := range assets {
			log.Printf("\tchecking %s...\n", asset.ID)
			a, err := assetDB.GetByID(asset.ID)
			if err != nil {
				log.Fatal(err)
			}
			if a == nil {
				log.Printf("\tinserting %s...\n", asset.ID)
				a, err = assetDB.Insert(asset)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Printf("\t%s already inserted\n", asset.ID)
			}
		}

		log.Println("done")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
