package main

import (
	"fmt"
	"os"

	"github.com/bitfield/uptimerobot"
)

func main() {
	apiKey := os.Args[1]
	utr := uptimerobot.New(apiKey)
	fmt.Println(utr.GetAccountDetails())
}
