package main

import (
	"bartender"
	"log"

	"github.com/abatewongc/bartender-bastion/client/league"
)

func main() {

	client, err := league.NewFromExisting()
	if err != nil {
		log.Fatalf("startup error: %v", err)
	}

	app := bartender.New(client)
	app.Listen()
}
