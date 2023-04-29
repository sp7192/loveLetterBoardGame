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
		players: []Player{
			{ID: 1},
			{ID: 2},
		},
		deck: deck.NewDeck(cards),
	}

	firstCount := g.deck.Count()

	// Prepare the game phase
	err := g.PreparePhase()

	// Check that there were no errors
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedCounts := firstCount - len(g.players)
	// Check that the deck was shuffled
	if g.deck.Count() != expectedCounts {
		t.Errorf("expected %d cards in deck, got %d", expectedCounts, g.deck.Count())
	}

	// Check that each player has one card in their hand
	for _, p := range g.players {
		if len(p.hand.cards) != 1 {
			t.Errorf("expected player %d to have 1 card, got %d", p.ID, len(p.hand.cards))
		}
	}
}
