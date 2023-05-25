package gameloop

import (
	"encoding/json"
	"loveLetterBoardGame/internals/configs"
	gamelogic "loveLetterBoardGame/internals/gamelogic"
	"loveLetterBoardGame/internals/server"
	"time"
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
	err = g.server.SendToAll(state)
	if err != nil {
		return nil
	}

	err = g.sendPlayerCardsInHandToAll()
	if err != nil {
		return err
	}

	err = g.runTurns()
	if err != nil {
		return err
	}
	return nil
}

func (g *GameLoop) sendPlayerCardsInHand(id uint) error {
	cards := g.gameLogic.GetPlayersCardsInHand(id)
	data, err := json.Marshal(cards)
	if err != nil {
		return err
	}
	g.server.SendTo(id, string(data))
	return nil
}

func (g *GameLoop) sendPlayerCardsInHandToAll() error {
	for _, p := range g.gameLogic.Players {
		err := g.sendPlayerCardsInHand(p.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GameLoop) sendGameStateToAll() error {
	state, err := g.gameLogic.GetGameState()
	if err != nil {
		return err
	}
	return g.server.SendToAll(state)
}

func (g *GameLoop) isGameEnded() bool {
	// TODO: To be implemeneted
	return false
}

func (g *GameLoop) runTurns() error {
	for {
		// 1. Draw phase
		ok := g.gameLogic.DrawPhase()
		if !ok {
			break
		}

		// 2. Send turn player card.
		err := g.sendPlayerCardsInHand(g.gameLogic.Players[g.gameLogic.PlayingPlayerIndex].ID)
		if err != nil {
			return err
		}

		// 3. Send game state to others.
		err = g.sendGameStateToAll()
		if err != nil {
			return err
		}

		// 4. Receive player action. (Random action if Timeout).
		msg, err := g.server.GetClientMessage()
		if err != nil {
			// handle timeout
			return err
		}

		// 5. Update the game based on player action
		g.gameLogic.UpdateGame(msg)

		// 6. Send game state to others.

		// 7. Check for game end condition.
		if g.isGameEnded() {
			break
		}

		// 8.Change to next playing player.

		// TODO Remove
		time.Sleep(5 * time.Second)
	}
	// 9. Find winner.
	return nil
}
