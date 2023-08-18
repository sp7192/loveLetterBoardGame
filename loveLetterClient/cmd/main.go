package main

import (
	"log"
	"loveLetterClient/internals/client"
	"loveLetterClient/internals/configs"
)

func main() {
	logger := log.Default()

	conf, err := configs.LoadConfigs("../internals/configs", "configs", "env")
	if err != nil {
		log.Fatal(err.Error())
	}

	logger.Printf("Client id is : %d\n", conf.ClientID)
	cl := client.NewClient(&conf, logger)
	cl.Run()
}
