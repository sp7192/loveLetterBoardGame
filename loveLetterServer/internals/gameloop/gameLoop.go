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

func (g *GameLoop) BeginGame() {
	g.gameLogic.PreparePhase()
	g.server.SendToAll(g.gameLogic)
	g.gameLogic.BeginTurns()
}
