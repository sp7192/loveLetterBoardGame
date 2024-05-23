package main

import (
	"log"
	"loveLetterBoardGame/cmd/client/internal/client"
	"loveLetterBoardGame/internal/configs"
)

func main() {
	logger := log.Default()

	conf, err := configs.LoadConfigs("../../configs", "configs", "env")
	if err != nil {
		log.Fatal(err.Error())
	}

	logger.Printf("Client id is : %d\n", conf.ClientID)
	cl := client.NewClient(&conf, logger)
	cl.Run()
}
