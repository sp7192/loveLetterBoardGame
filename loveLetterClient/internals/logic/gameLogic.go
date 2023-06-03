package logic

import "loveLetterClient/internals/models"

type GameLogic struct {
	ownHand         models.Hand
	playersIdInGame []uint
	playedCards     []models.Card
}

func (g *GameLogic) ParseMessage(msg string) error {
	return nil
}
