package deck

import (
	"loveLetterBoardGame/internals/gamelogic/card"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDeck_Draw(t *testing.T) {
	cards := []card.Card{{Number: 1}, {Number: 2}, {Number: 3}}
	deck := NewDeck(cards)

	require.NotEmpty(t, deck)

	// Draw all cards from the deck
	for i := 0; i < len(cards); i++ {
		card, ok := deck.Draw()
		require.True(t, ok)
		require.Equal(t, cards[deck.count], card)
	}

	// Try to draw one more card from an empty deck
	c, ok := deck.Draw()
	require.False(t, ok)
	require.Equal(t, card.Card{}, c)
}

func TestShuffle(t *testing.T) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a deck with known cards
	cards := []card.Card{
		{1, nil},
		{2, nil},
		{3, nil},
		{4, nil},
		{5, nil},
	}
	oldCards := make([]card.Card, len(cards))
	copy(oldCards, cards)
	deck := NewDeck(cards)

	// Shuffle the deck 10 times
	for i := 0; i < 10; i++ {
		deck.Shuffle()
		if !reflect.DeepEqual(deck.cards, oldCards) {
			// The order of the cards has changed, so no need to continue shuffling
			return
		}
	}

	// The order of the cards has not changed after 10 shuffles, report error
	t.Errorf("Shuffle() failed. The order of cards has not changed after 10 shuffles")
}
