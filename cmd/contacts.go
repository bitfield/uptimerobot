package cmd

import (
	"fmt"
	"log"

	uptimerobot "github.com/bitfield/uptimerobot/pkg"
	"github.com/spf13/cobra"
)

var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "list alert contacts",
	Long:  `Show all alert contacts associated with the account`,
	Run: func(cmd *cobra.Command, args []string) {
		utr := uptimerobot.New(apiKey)
		contacts, err := utr.GetAlertContacts()
		if err != nil {
			log.Fatal(err)
		}
		if len(contacts) == 0 {
			fmt.Println("No contacts found")
		}
		for _, c := range contacts {
			fmt.Println(c)
			fmt.Println()
		}
	},
}

func init() {
	RootCmd.AddCommand(contactsCmd)
}
