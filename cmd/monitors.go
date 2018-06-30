package cmd

import (
	"fmt"
	"log"

	uptimerobot "github.com/bitfield/uptimerobot/pkg"
	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitors",
	Short: "lists monitors",
	Long:  `Lists all monitors associated with the account`,
	Run: func(cmd *cobra.Command, args []string) {
		utr := uptimerobot.New(apiKey)
		monitors, err := utr.GetMonitors()
		if err != nil {
			log.Fatal(err)
		}
		if len(monitors) == 0 {
			log.Fatal("No matching monitors found")
		}
		for _, m := range monitors {
			fmt.Println(m)
			fmt.Println()
		}
	},
}

func init() {
	RootCmd.AddCommand(monitorCmd)
}
