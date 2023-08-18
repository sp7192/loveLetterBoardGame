package logic

import (
	"encoding/json"
	"fmt"
	"loveLetterClient/internals/models"
)

type GameLogic struct {
	OwnHand          models.Hand
	playersIdInGame  []uint
	playedCards      []models.Card
	playingPlayerId  uint
	SendMessageQueue chan string
}

func NewGameLogic() *GameLogic {
	return &GameLogic{
		OwnHand: models.Hand{
			Cards: make([]models.Card, 0, 2),
		},
		playersIdInGame:  make([]uint, 0, 10),
		playedCards:      make([]models.Card, 0, 32),
		playingPlayerId:  0,
		SendMessageQueue: make(chan string),
	}
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

func (g *GameLogic) PlayTurn() error {
	// TODO : to get from player
	action := models.ClientAction{
		PlayedCardNumber: g.OwnHand.Cards[0].Number,
	}
	str, err := json.Marshal(action)
	if err != nil {
		return err
	}
	g.SendMessageQueue <- string(str)
	return nil
}

func (g *GameLogic) update(msg models.Message) error {
	switch msg.Type {
	case models.InitDrawMessage:
		fmt.Printf(">> Init Draw message, Data : %s\n\n", msg.Payload)
		err := json.Unmarshal([]byte(msg.Payload), &g.OwnHand.Cards)
		if err != nil {
			return fmt.Errorf("error in initDraw message : %s\n", err.Error())
		}
		return nil
	case models.TurnDrawMessage:
		// TODO : to be completed
		fmt.Printf(">> Turn Draw message, Data : %s\n\n", msg.Payload)
		return g.PlayTurn()
	case models.InfoMessage:
		// TODO : to be completed
		fmt.Printf(">> Info message, Data : %s\n\n", msg.Payload)
	case models.UpdateMessage:
		// TODO : to be completed
		fmt.Printf(">> Update message, Data : %s\n\n", msg.Payload)
	case models.PlayedMessage:
		// TODO : to be completed
		fmt.Printf(">> Played message, Data : %s\n\n", msg.Payload)
	default:
		return fmt.Errorf("Message type %s, not supported", msg.Type)
	}
	return nil
}
