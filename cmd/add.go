package cmd

import (
	"fmt"
	"log"

	uptimerobot "github.com/bitfield/uptimerobot/pkg"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a new monitor",
	Long:  `Create a new monitor with the specified name and URL`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		utr := uptimerobot.New(apiKey)
		m := uptimerobot.Monitor{
			URL:          args[0],
			FriendlyName: args[1],
			Type:         uptimerobot.MonitorHTTP,
		}
		new, err := utr.NewMonitor(m)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("New monitor created with ID %d\n", new.ID)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
