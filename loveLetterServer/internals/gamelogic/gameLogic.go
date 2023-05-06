package gamelogic

import (
	"fmt"
	"loveLetterBoardGame/internals/gamelogic/card"
	"loveLetterBoardGame/internals/gamelogic/deck"
	"time"
)

type GameLogic struct {
	Players []Player  `json:"players"`
	Deck    deck.Deck `json:"deck"`
}

func NewGameLogic(mode string, players []Player) GameLogic {
	cards := card.NewCardsSet("TEST")
	return GameLogic{
		Deck:    deck.NewDeck(cards),
		Players: players,
	}
}

func (g *GameLogic) PreparePhase() error {
	g.Deck.Shuffle()
	for i := range g.Players {
		card, ok := g.Deck.Draw()
		if !ok {
			return fmt.Errorf("not enough cards in deck")
		}
		g.Players[i].hand.cards = append(g.Players[i].hand.cards, card)
	}
	return nil
}

func (g *GameLogic) BeginTurns() {
	for {
		fmt.Println("SIMULATE TURNS.. TODO REMOVE.")
		time.Sleep(time.Second)
	}
}

func (g *GameLogic) DrawPhase() {

}
