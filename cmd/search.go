package cmd

import (
	"fmt"
	"log"

	uptimerobot "github.com/bitfield/uptimerobot/pkg"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search monitors",
	Long:  `Lists all monitors matching a search string`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utr := uptimerobot.New(apiKey)
		monitors, err := utr.GetMonitorsBySearch(args[0])
		if err != nil {
			log.Fatal(err)
		}
		if len(monitors) == 0 {
			log.Fatal("No matching monitors found")
		}
		for _, m := range monitors {
			fmt.Println(m)
		}
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)
}
