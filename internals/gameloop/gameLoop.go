package gameloop

import (
	"encoding/json"
	"log"
	"loveLetterBoardGame/internals/configs"
	gamelogic "loveLetterBoardGame/internals/gamelogic"
	"loveLetterBoardGame/internals/server"
	"loveLetterBoardGame/models"
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
	err = g.server.SendToAllWithAck(state)
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

func (g *GameLoop) sendPlayerCardsInHand(id uint, msgType models.MessageType, hasAck bool) error {
	cards := g.gameLogic.GetPlayersCardsInHand(id)
	data, err := json.MarshalIndent(cards, "", "	")
	if err != nil {
		return err
	}
	// TODO: Refactor
	if !hasAck {
		err = g.server.SendTo(id, msgType, string(data))
	} else {
		err = g.server.SendAndReceiveAck(id, msgType, string(data))
	}
	if err != nil {
		return err
	}
	return nil
}

func (g *GameLoop) sendPlayerCardsInHandToAll() error {
	for _, p := range g.gameLogic.Players {
		err := g.sendPlayerCardsInHand(p.ID, models.InitDrawMessage, true)
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
	return g.server.SendToAllWithAck(state)
}

func (g *GameLoop) isGameEnded() bool {
	// TODO: To be implemeneted
	return false
}

func (g *GameLoop) runTurns() error {
	turn := 0
	for {
		turn++
		g.logger.Println("\n==============================\n1. Draw phase, turn = ", turn)
		ok := g.gameLogic.DrawPhase()
		if !ok {
			break
		}

		g.logger.Println("2. Send game state to others")
		err := g.sendGameStateToAll()
		if err != nil {
			return err
		}

		g.logger.Println("3. Send turn player card")
		err = g.sendPlayerCardsInHand(g.gameLogic.Players[g.gameLogic.PlayingPlayerIndex].ID, models.TurnDrawMessage, false)
		if err != nil {
			return err
		}

		g.logger.Println("4. Receive player action. (Random action if Timeout).")
		msg, err := g.server.GetClientMessage()
		if err != nil {
			return err
		}
		g.logger.Printf("Client Message: %v\n", msg)

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

		g.logger.Print("8. Change to next playing player.\n\n")
		g.gameLogic.ChangePlayingPlayer()
	}
	// 9. Find winner.
	return nil
}
