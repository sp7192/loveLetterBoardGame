package deck

import (
	"loveLetterBoardGame/internals/gamelogic/card"
	"testing"

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
