package gamelogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePlayersFromIDs(t *testing.T) {
	ids := []uint{1, 2, 3}
	players := CreatePlayersFromIDs(ids)

	assert.Len(t, players, len(ids))
	for i, player := range players {
		assert.Equal(t, ids[i], player.ID)
		assert.Equal(t, uint(0), player.totalScore)
		assert.Equal(t, Hand{}, player.Hand)
		assert.True(t, player.isInThisRound)
	}
}
