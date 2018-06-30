package cmd

import (
	"fmt"
	"log"

	uptimerobot "github.com/bitfield/uptimerobot/pkg"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "get account details",
	Long:  `Show the account details associated with the account`,
	Run: func(cmd *cobra.Command, args []string) {
		utr := uptimerobot.New(apiKey)
		account, err := utr.GetAccountDetails()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(account)
	},
}

func init() {
	RootCmd.AddCommand(accountCmd)
}
