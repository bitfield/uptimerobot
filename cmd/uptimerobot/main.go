package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bitfield/uptimerobot"
)

func main() {
	apiKey := os.Args[1]
	utr := uptimerobot.New(apiKey)
	monitors, err := utr.GetMonitors()
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range monitors {
		fmt.Println(m.FriendlyName)
	}
}
