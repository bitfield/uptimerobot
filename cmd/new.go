package cmd

import (
	"fmt"
	"log"

	uptimerobot "github.com/bitfield/uptimerobot/pkg"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "add a new monitor",
	Long:  `Create a new monitor with the specified URL and friendly name`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		m := uptimerobot.Monitor{
			URL:           args[0],
			FriendlyName:  args[1],
			Type:          uptimerobot.MonitorType("HTTP"),
			AlertContacts: contacts,
		}
		new, err := client.NewMonitor(m)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("New monitor created with ID %d\n", new.ID)
	},
}

var contacts []string

func init() {
	newCmd.Flags().StringSliceVarP(&contacts, "contacts", "c", []string{}, "Comma-separated list of contact IDs to notify")
	RootCmd.AddCommand(newCmd)
}
