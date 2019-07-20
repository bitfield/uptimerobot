package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitors",
	Short: "lists monitors",
	Long:  `Lists all monitors associated with the account`,
	Run: func(cmd *cobra.Command, args []string) {
		result := client.GetMonitorsChan()

		done := make(chan struct{})
		go func() {
			for m := range result.Monitors {
				fmt.Println(m)
				fmt.Println()
			}
			close(done)
		}()
		select {
		case <-done:
			return
		case err := <-result.Error:
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(monitorCmd)
}
