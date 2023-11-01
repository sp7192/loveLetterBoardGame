package logic

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"loveLetterClient/internals/models"
)

type GameLogic struct {
	OwnHand          models.Hand
	playersIdInGame  []uint
	playedCards      []models.Card
	playingPlayerId  uint
	SendMessageQueue chan string
	logger           *log.Logger
	inputScanner     *bufio.Scanner
}

func NewGameLogic(l *log.Logger, sc *bufio.Scanner) *GameLogic {
	return &GameLogic{
		OwnHand: models.Hand{
			Cards: make([]models.Card, 0, 2),
		},
		playersIdInGame:  make([]uint, 0, 10),
		playedCards:      make([]models.Card, 0, 32),
		playingPlayerId:  0,
		SendMessageQueue: make(chan string),
		logger:           l,
		inputScanner:     sc,
	}
}

func (g *GameLogic) getUserInput() uint {
	g.logger.Printf("Choose your card:\nEnter 0 for Card:%v\nEnter 1 for Card:%v\n\n", g.OwnHand.Cards[0], g.OwnHand.Cards[1])

	for g.inputScanner.Scan() {
		text := g.inputScanner.Text()
		if text == "0" {
			return 0
		} else if text == "1" {
			return 1
		} else {
			g.logger.Printf("Wrong Input!\n")
			g.logger.Printf("Choose your card:\nEnter 0 for Card:%v\nEnter 1 for Card:%v\n\n", g.OwnHand.Cards[0], g.OwnHand.Cards[1])
		}
	}
	return 0
}

func (g *GameLogic) SendReceiveAck(msg models.Message) error {

	sendMessage := models.Message{
		Payload: msg.Type,
		Type:    models.AckMessage,
	}

	str, err := json.MarshalIndent(sendMessage, "", "	")
	if err != nil {
		return err
	}
	// TODO: Add timeout
	g.SendMessageQueue <- string(str)
	return nil
}

func (g *GameLogic) Update(msg models.Message) error {
	err := g.update(msg)
	if err != nil {
		return err
	}
	return nil
}

func (g *GameLogic) PlayTurn() error {
	// TODO : to add timeout.
	// TODO : to add target player.
	// TODO : to add some action(like guessing the card).
	index := g.getUserInput()
	action := models.ClientAction{
		PlayedCardNumber: g.OwnHand.Cards[index].Number,
	}
	str, err := json.MarshalIndent(action, "", "	")
	if err != nil {
		return err
	}
	// TODO: Add timeout
	g.SendMessageQueue <- string(str)
	return nil
}

func (g *GameLogic) update(msg models.Message) error {
	switch msg.Type {
	case models.InitDrawMessage:
		g.logger.Printf(">> Init Draw message, Data : %s\n\n", msg.Payload)
		err := json.Unmarshal([]byte(msg.Payload), &g.OwnHand.Cards)
		if err != nil {
			return fmt.Errorf("error in initDraw message : %s\n", err.Error())
		}
		return nil
	case models.TurnDrawMessage:
		// TODO : to be completed
		g.logger.Printf(">> Turn Draw message, Data : %s\n\n", msg.Payload)
		err := json.Unmarshal([]byte(msg.Payload), &g.OwnHand.Cards)
		if err != nil {
			return fmt.Errorf("error in initDraw message : %s\n", err.Error())
		}
		return g.PlayTurn()
	case models.InfoMessage:
		// TODO : to be completed
		g.logger.Printf(">> Info message, Data : %s\n\n", msg.Payload)
	case models.UpdateMessage:
		// TODO : to be completed
		g.logger.Printf(">> Update message, Data : %s\n\n", msg.Payload)
	case models.PlayedMessage:
		// TODO : to be completed
		g.logger.Printf(">> Played message, Data : %s\n\n", msg.Payload)
	default:
		return fmt.Errorf("Message type %s, not supported", msg.Type)
	}
	return nil
}
