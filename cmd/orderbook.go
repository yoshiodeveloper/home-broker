package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// orderbookCmd represents the orderbook command
var orderbookCmd = &cobra.Command{
	Use:   "orderbook",
	Short: "Starts the order book service",
	Long:  `You must have only one instance of the order book service.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("orderbook called")
	},
}

func init() {
	rootCmd.AddCommand(orderbookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// orderbookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// orderbookCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
