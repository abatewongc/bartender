package main

import (
	"bartender"
	"fmt"
	"os"
	"time"

	"github.com/coltiebaby/bastion/client/league"
)

func main() {

	client, err := league.NewFromExisting()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	app := bartender.New(bartender.Config{
		SkinBlacklist:  map[float64]struct{}{},
		Tickrate:       time.Millisecond * 500,
		InGameTickrate: time.Minute * 8,
	}, client)
	app.Listen()
}
