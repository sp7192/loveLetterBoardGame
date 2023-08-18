package gameloop

import (
	"encoding/json"
	"log"
	"loveLetterBoardGame/internals/configs"
	gamelogic "loveLetterBoardGame/internals/gamelogic"
	"loveLetterBoardGame/internals/server"
	"loveLetterBoardGame/models"
	"time"
)

type GameLoop struct {
	server    *server.Server
	gameLogic *gamelogic.GameLogic
	configs   *configs.Configs
	logger    *log.Logger
}

func NewGameLoop(s *server.Server, g *gamelogic.GameLogic, c *configs.Configs, l *log.Logger) GameLoop {
	return GameLoop{server: s, gameLogic: g, configs: c, logger: l}
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

func (g *GameLoop) sendPlayerCardsInHand(id uint, msgType models.MessageType) error {
	cards := g.gameLogic.GetPlayersCardsInHand(id)
	data, err := json.Marshal(cards)
	if err != nil {
		return err
	}
	g.server.SendTo(id, msgType, string(data))
	return nil
}

func (g *GameLoop) sendPlayerCardsInHandToAll() error {
	for _, p := range g.gameLogic.Players {
		err := g.sendPlayerCardsInHand(p.ID, models.InitDrawMessage)
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
		g.logger.Println("1. Draw phase")
		ok := g.gameLogic.DrawPhase()
		if !ok {
			break
		}

		g.logger.Println("2. Send turn player card")
		err := g.sendPlayerCardsInHand(g.gameLogic.Players[g.gameLogic.PlayingPlayerIndex].ID, models.TurnDrawMessage)
		if err != nil {
			return err
		}

		g.logger.Println("3. Send game state to others")
		err = g.sendGameStateToAll()
		if err != nil {
			return err
		}

		g.logger.Println("4. Receive player action. (Random action if Timeout).")
		msg, err := g.server.GetClientMessage()
		if err != nil {
			return err
		}

		g.logger.Println("5. Update the game based on player action")
		g.gameLogic.UpdateGame(msg)

		g.logger.Println("6. Send game state to others.")
		err = g.sendGameStateToAll()
		if err != nil {
			return err
		}

		g.logger.Println("7. Check for game end condition.")
		if g.isGameEnded() {
			break
		}

		g.logger.Println("8. Change to next playing player.")

		// TODO Remove
		time.Sleep(5 * time.Second)
	}
	// 9. Find winner.
	return nil
}
