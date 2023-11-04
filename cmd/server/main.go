package main

import (
	"log"
	"loveLetterBoardGame/internals/configs"
	"loveLetterBoardGame/internals/gamelogic"
	"loveLetterBoardGame/internals/gameloop"
	"loveLetterBoardGame/internals/server"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	logger := log.Default()
	conf, err := configs.LoadConfigs("configs", "configs", "env")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create a new Server
	srv := server.NewServer(conf, logger)

	err = srv.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	players := gamelogic.CreatePlayersFromIDs(srv.GetClientsIds())
	logic := gamelogic.NewGameLogic("", players, logger)

	loop := gameloop.NewGameLoop(&srv, &logic, &conf, logger)
	_ = loop.BeginGame()
}
