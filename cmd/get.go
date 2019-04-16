package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get monitor by ID",
	Long:  `Show the monitor details for the specified monitor ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		monitor, err := client.GetMonitorByID(ID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(monitor)
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
