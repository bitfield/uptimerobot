package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/bitfield/uptimerobot/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "uptimerobot",
	Short: "uptimerobot is a client for the UptimeRobot V2 API",
	Long: `uptimerobot is a command-line client for the UptimeRobot monitoring
service. It allows you to search for existing monitors, delete monitors,
and create new monitors. You can also inspect your account details and
any alert contacts you have configured.

For more information, see https://github.com/bitfield/uptimerobot`,
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
var client uptimerobot.Client

func init() {
	viper.SetConfigName(".uptimerobot")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("failed to read config: %v\n", err)
		}
	}
	viper.SetEnvPrefix("uptimerobot")
	viper.AutomaticEnv()
	cobra.OnInitialize(func() {
		client = uptimerobot.New(viper.GetString("apiKey"))
		client.Debug = debug
	})
	RootCmd.PersistentFlags().StringVar(&apiKey, "apiKey", "", "UptimeRobot API key")
	viper.BindPFlag("apiKey", RootCmd.PersistentFlags().Lookup("apiKey"))
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug mode (show API request without making it)")
}
