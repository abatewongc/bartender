package main

import (
	"bartender"
	"fmt"
	"os"

	"github.com/abatewongc/bartender-bastion/client/league"
)

func main() {

	client, err := league.NewFromExisting()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	app := bartender.New(client)
	app.Listen()
}
