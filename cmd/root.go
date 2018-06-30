package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "uptimerobot",
	Short: "uptimerobot is a client for the UptimeRobot V2 API",
	Long:  `Insert helpful text here`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var apiKey string
var debug bool

func init() {
	RootCmd.PersistentFlags().StringVar(&apiKey, "apiKey", "", "UptimeRobot API key")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug mode (show API request without making it)")
}
