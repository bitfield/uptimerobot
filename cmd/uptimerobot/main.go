package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bitfield/uptimerobot"
)

func main() {
	var apiKey = flag.String("api-key", "", "UptimeRobot API key")
	var debug = flag.Bool("debug", false, "Debug mode (show API request without making it)")
	var get = flag.String("get", "", "Return all monitors matching string (or 'all')")
	var new = flag.String("new", "", "Create a new monitor with specified name (requires '-url')")
	var url = flag.String("url", "", "URL for new monitor check (used with '-new')")
	flag.Parse()
	if *apiKey == "" {
		usageError("Please specify an UptimeRobot API key to use")
	}
	utr := uptimerobot.New(*apiKey)
	utr.Debug = *debug
	if *get != "" {
		var monitors []uptimerobot.Monitor
		var err error
		if *get == "all" {
			monitors, err = utr.GetMonitors()
		} else {
			monitors, err = utr.GetMonitorsBySearch(*get)
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(monitors) == 0 {
			log.Fatal("No matching monitors found")
		}
		for _, m := range monitors {
			fmt.Println(m)
		}
		os.Exit(0)
	}
	if *new != "" {
		if *url == "" {
			usageError("Please specify a URL for the new monitor")
		}
		m := uptimerobot.Monitor{
			FriendlyName: *new,
			URL:          *url,
			Type:         uptimerobot.MonitorHTTP,
		}
		result, err := utr.NewMonitor(m)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("New monitor created with ID %d\n", result.ID)
		os.Exit(0)
	}
}

func usageError(msg string) {
	log.Println(msg)
	flag.Usage()
	os.Exit(1)
}
