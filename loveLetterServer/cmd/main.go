package main

import (
	"log"
	"loveLetterBoardGame/internals/configs"
	"loveLetterBoardGame/internals/gamelogic"
	"loveLetterBoardGame/internals/gameloop"
	"loveLetterBoardGame/internals/server"
)

func main() {
	conf, err := configs.LoadConfigs("../internals/configs", "configs", "env")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create a new Server
	srv := server.NewServer(conf)

	err = srv.Start()
	if err != nil {
		log.Fatal(err.Error())
	}


	logic := gamelogic.NewGameLogic("", []gamelogic.Player{})

	loop := gameloop.NewGameLoop(&srv, &logic, &conf)
	loop.BeginTurn()
}
