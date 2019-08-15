package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search monitors",
	Long:  `Lists all monitors matching a search string`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		monitors, err := client.SearchMonitors(args[0])
		if err != nil {
			log.Fatal(err)
		}
		if len(monitors) == 0 {
			fmt.Println("No matching monitors found")
			os.Exit(1)
		}
		for _, m := range monitors {
			fmt.Println(m)
			fmt.Println()
		}
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)
}
