package gamelogic

import (
	"fmt"
	"loveLetterBoardGame/internals/gamelogic/card"
	"loveLetterBoardGame/internals/gamelogic/deck"
	"math/rand"
	"time"
)

type GameLogic struct {
	Players         []Player
	Deck            deck.Deck
	PlayingPlayerId uint
}

func NewGameLogic(mode string, players []Player) GameLogic {
	cards := card.NewCardsSet("TEST")
	return GameLogic{
		Deck:    deck.NewDeck(cards),
		Players: players,
	}
}

func (g *GameLogic) getStartingPlayerId() uint {
	return g.Players[rand.Intn(len(g.Players))].ID
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
	g.PlayingPlayerId = g.getStartingPlayerId()
	return nil
}

func (g *GameLogic) isPlayerIdValid(playerId uint) bool {
	for _, p := range g.Players {
		if playerId == p.ID {
			return true
		}
	}
	return false
}

func (g *GameLogic) GetGameState() (GameState, error) {
	var ret GameState

	ret.DeckCardsCount = uint(g.Deck.Count())
	for _, p := range g.Players {
		if p.isInThisRound {
			ret.PlayersIdInGame = append(ret.PlayersIdInGame, p.ID)
		}
	}

	ret.PlayingPlayerId = g.PlayingPlayerId

	return ret, nil
}

func (g *GameLogic) BeginTurns() {
	for {
		fmt.Println("SIMULATE TURNS.. TODO REMOVE.")
		time.Sleep(time.Second)
	}
}

func (g *GameLogic) DrawPhase() {

}
