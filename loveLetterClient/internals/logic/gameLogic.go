package logic

import "loveLetterClient/internals/models"

type GameLogic struct {
	ownHand         models.Hand
	playersIdInGame []uint
	playedCards     []models.Card
}
