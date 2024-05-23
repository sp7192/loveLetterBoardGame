package main

import (
	"log"
	"loveLetterBoardGame/cmd/server/internal/gamelogic"
	"loveLetterBoardGame/cmd/server/internal/gameloop"
	"loveLetterBoardGame/cmd/server/internal/server"
	"loveLetterBoardGame/internal/configs"
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
	err = loop.BeginGame()
	if err != nil {
		log.Fatal(err.Error())
	}
}
