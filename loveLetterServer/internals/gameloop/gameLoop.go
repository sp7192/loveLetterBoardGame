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

	err = g.runTurns()
	if err != nil {
		return err
	}
	return nil
}

func (g *GameLoop) isGameEnded() bool {
	// TODO: To be implemeneted
	return false
}

func (g *GameLoop) runTurns() error {
	for {
		// Send turn player card.
		// send game state to others.
		// receive player action. (Random action if Timeout).
		// send game state to others.
		// check for game end condition.
		if g.isGameEnded() {
			break
		}
		// change to next playing player.
	}
	return nil
}
