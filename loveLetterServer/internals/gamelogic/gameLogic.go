package gamelogic

import (
	"fmt"
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

func (g *GameLogic) PreparePhase() error {
	g.deck.Shuffle()
	for i := range g.players {
		card, ok := g.deck.Draw()
		if !ok {
			return fmt.Errorf("not enough cards in deck")
		}
		g.players[i].hand.cards = append(g.players[i].hand.cards, card)
	}
	return nil
}

func (g *GameLogic) BeginTurns() {
}

func (g *GameLogic) DrawPhase() {

}
