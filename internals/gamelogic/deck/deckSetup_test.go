package deck

import (
	"loveLetterBoardGame/internals/gamelogic/card"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDeckSetups(t *testing.T) {
	deckSetups, err := ParseDeckSetups("decks_setup.yaml")
	assert.NoError(t, err, "Error parsing deck setups")

	assert.NotEmpty(t, deckSetups, "Deck setups should not be empty")

	for _, setup := range deckSetups {
		assert.NotEmpty(t, setup.Type, "Deck setup type should not be empty")

		// Ensure each card in the setup is valid
		for _, c := range setup.Cards {
			assert.NotEmpty(t, c.Name, "Card name should not be empty")
			assert.NotEqual(t, 0, c.Number, "Card number should not be zero")
			// Add more assertions based on your card structure
		}
	}
}

func TestFindByType(t *testing.T) {
	deckSetups := DeckSetups{
		{
			Type:  "normal",
			Cards: []card.Card{
				// Add example cards here
			},
		},
		{
			Type:  "extended",
			Cards: []card.Card{
				// Add example cards here
			},
		},
	}

	setup, err := deckSetups.FindByType(Normal)
	assert.NoError(t, err, "Error finding deck setup by type")
	assert.Equal(t, "normal", setup.Type, "Found setup type should be 'normal'")

	setup, err = deckSetups.FindByType(Extended)
	assert.NoError(t, err, "Error finding deck setup by type")
	assert.Equal(t, "extended", setup.Type, "Found setup type should be 'extended'")

	// Test with a non-existent type
	setup, err = deckSetups.FindByType(DeckType("nonexistent"))
	assert.Error(t, err, "Expected error for non-existent type")
	assert.Equal(t, DeckSetup{}, setup, "Found setup should be empty for non-existent type")
}
