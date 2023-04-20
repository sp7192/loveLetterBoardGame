package gamemanager

import (
	gamelogic "loveLetterBoardGame/internals/gameLogic"
	"loveLetterBoardGame/internals/server"
)

type GameManager struct {
	server    *server.Server
	gameLogic *gamelogic.GameLogic
}

func NewGameManager(s *server.Server, g *gamelogic.GameLogic) GameManager {
	return GameManager{server: s, gameLogic: g}
}

func (gm *GameManager) RunGame() error {
	err := gm.server.Start()
	if err != nil {
		return err
	}
	
	return nil
}
