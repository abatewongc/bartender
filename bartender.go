package main

import (
	"bartender/internal/service"
	"fmt"
	"github.com/coltiebaby/bastion/client/league"
	"os"
)

func main() {

	client, err := league.NewFromExisting()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	bartenderService := service.BuildBartenderService(client)
	bartenderService.Listen()

	return
}
