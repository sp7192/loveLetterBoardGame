package gameloop

import (
	"loveLetterBoardGame/internals/configs"
	gamelogic "loveLetterBoardGame/internals/gamelogic"
	"loveLetterBoardGame/internals/server"
)

type GameLoop struct {
	server    *server.Server
	gameLogic *gamelogic.GameLogic
	configs   *configs.Configs
}

func NewGameLoop(s *server.Server, g *gamelogic.GameLogic, c *configs.Configs) GameLoop {
	return GameLoop{server: s, gameLogic: g, configs: c}
}

func (g *GameLoop) RunGame() error {
	err := g.server.Start()
	if err != nil {
		return err
	}

	return nil
}

func (g *GameLoop) BeginTurn() {
}
