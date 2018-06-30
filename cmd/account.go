package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "get account details",
	Long:  `Show the account details associated with the API key.`,
	Run: func(cmd *cobra.Command, args []string) {
		account, err := client.GetAccountDetails()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(account)
	},
}

func init() {
	RootCmd.AddCommand(accountCmd)
}
