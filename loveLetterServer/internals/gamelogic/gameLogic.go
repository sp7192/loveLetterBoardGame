package gamelogic

import (
	"fmt"
	"loveLetterBoardGame/internals/gamelogic/card"
	"loveLetterBoardGame/internals/gamelogic/deck"
	"time"
)

type GameLogic struct {
	Players []Player
	Deck    deck.Deck
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

func (g *GameLogic) isPlayerIdValid(playerId uint) bool {
	for _, p := range g.Players {
		if playerId == p.ID {
			return true
		}
	}
	return false
}

func (g *GameLogic) GetGameState(playerId uint) (GameState, error) {
	var ret GameState

	if !g.isPlayerIdValid(playerId) {
		return ret, fmt.Errorf("Player id : %d, is invalid", playerId)
	}

	ret.DeckCardsCount = uint(g.Deck.Count())
	for _, p := range g.Players {
		if p.isInThisRound {
			ret.PlayersIdInGame = append(ret.PlayersIdInGame, p.ID)
		}
	}

	ret.PlayingPlayerId = playerId

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
