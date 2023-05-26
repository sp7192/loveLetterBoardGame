package gamelogic

import (
	"loveLetterBoardGame/internals/gamelogic/card"
	"loveLetterBoardGame/internals/gamelogic/deck"
	"testing"
)

func TestGameLogic_PreparePhase(t *testing.T) {
	// Create a game logic instance with two players and a deck with two cards
	cards := card.NewCardsSet("TEST")
	g := &GameLogic{
		Players: []Player{
			{ID: 1},
			{ID: 2},
		},
		Deck: deck.NewDeck(cards),
	}

	expectedCounts := g.Deck.Count() - len(g.Players)

	// Prepare the game phase
	err := g.PreparePhase()

	// Check that there were no errors
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that the deck was shuffled
	if g.Deck.Count() != expectedCounts {
		t.Errorf("expected %d cards in deck, got %d", expectedCounts, g.Deck.Count())
	}

	// Check that each player has one card in their hand
	for _, p := range g.Players {
		if len(p.hand.cards) != 1 {
			t.Errorf("expected player %d to have 1 card, got %d", p.ID, len(p.hand.cards))
		}
	}
}
