package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.11.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  `Show uptimerobot client version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("uptimerobot version", version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
