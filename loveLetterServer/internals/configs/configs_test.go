package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigs(t *testing.T) {
	t.Run("load config from file", func(t *testing.T) {
		conf, err := LoadConfigs("./testdata", "testconfigs", "env")
		assert.NoError(t, err)
		assert.Equal(t, uint(4), conf.PlayersInRoomCount)
	})

	t.Run("load config from env", func(t *testing.T) {
		os.Setenv("PLAYERS_IN_ROOM_COUNT", "5")
		conf, err := LoadConfigs("./testdata", "testconfigs", "env")
		assert.NoError(t, err)
		assert.Equal(t, uint(5), conf.PlayersInRoomCount)
	})
}
