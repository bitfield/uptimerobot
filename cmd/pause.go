package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bitfield/uptimerobot/pkg"
	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "pause a monitor",
	Long:  `Pause the monitor with the specified ID`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		m := uptimerobot.Monitor{
			ID: ID,
		}
		new, err := client.PauseMonitor(m)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Monitor ID %d paused\n", new.ID)
	},
}

func init() {
	RootCmd.AddCommand(pauseCmd)
}
