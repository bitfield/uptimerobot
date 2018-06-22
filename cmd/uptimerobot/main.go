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
	var get = flag.String("get", "", "Return all monitors matching string (or 'all')")
	flag.Parse()
	if *apiKey == "" {
		usageError("Please specify an UptimeRobot API key to use")
	}
	utr := uptimerobot.New(*apiKey)
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
	}
}

func usageError(msg string) {
	log.Println(msg)
	flag.Usage()
	os.Exit(1)
}
