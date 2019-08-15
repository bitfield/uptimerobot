package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a monitor",
	Long:  `Delete the monitor with the specified ID`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		if err = client.DeleteMonitor(ID); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Monitor ID %d successfully deleted\n", ID)
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
