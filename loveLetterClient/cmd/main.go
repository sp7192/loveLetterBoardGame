package main

import (
	"fmt"
	"log"
	"loveLetterClient/internals/client"
	"loveLetterClient/internals/configs"
)

func main() {
	conf, err := configs.LoadConfigs("../internals/configs", "configs", "env")
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Client id is : %d\n", conf.ClientID)
	cl := client.NewClient(&conf)
	cl.Run()
}
