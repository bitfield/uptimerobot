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
	a, err := utr.GetAccountDetails()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)
}
