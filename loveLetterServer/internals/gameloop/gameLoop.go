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

func (g *GameLoop) BeginGame() error {
	g.gameLogic.PreparePhase()

	state, err := g.gameLogic.GetGameState()
	if err != nil {
		return err
	}
	g.server.SendToAll(state)

	g.gameLogic.BeginTurns()
	return nil
}
