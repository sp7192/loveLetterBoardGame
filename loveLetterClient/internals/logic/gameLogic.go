package logic

import (
	"encoding/json"
	"fmt"
	"loveLetterClient/internals/models"
)

type GameLogic struct {
	ownHand         models.Hand
	playersIdInGame []uint
	playedCards     []models.Card
	playingPlayerId uint
}

func (g *GameLogic) ParseMessage(strMsg string) error {
	var msg models.Message
	err := json.Unmarshal([]byte(strMsg), &msg)
	if err != nil {
		return err
	}
	err = g.update(msg)
	if err != nil {
		return err
	}
	return nil
}

func (g *GameLogic) update(msg models.Message) error {
	switch msg.Type {
	case models.DrawMessage:
		// TODO : to be completed
		fmt.Printf("Draw message, Data : %s\n", msg.Payload)
	case models.InfoMessage:
		// TODO : to be completed
		fmt.Printf("Info message, Data : %s\n", msg.Payload)
	case models.UpdateMessage:
		// TODO : to be completed
		fmt.Printf("Update message, Data : %s\n", msg.Payload)
	case models.PlayedMessage:
		// TODO : to be completed
		fmt.Printf("Played message, Data : %s\n", msg.Payload)
	default:
		return fmt.Errorf("Message type %s, not supported", msg.Type)
	}
	return nil
}
