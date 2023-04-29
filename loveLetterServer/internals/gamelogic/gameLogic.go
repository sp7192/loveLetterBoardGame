package gamelogic

import (
	"loveLetterBoardGame/internals/gamelogic/card"
	"loveLetterBoardGame/internals/gamelogic/deck"
)

type GameLogic struct {
	players []Player
	deck    deck.Deck
}

func NewGameLogic(mode string, players []Player) GameLogic {
	cards := card.NewCardsSet("TEST")
	return GameLogic{
		deck:    deck.NewDeck(cards),
		players: players,
	}
}

func (g *GameLogic) DrawPhase() {

}
