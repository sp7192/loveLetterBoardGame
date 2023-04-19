package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigs(t *testing.T) {
	conf, err := LoadConfigs("./testdata/", "testconfigs", "env")
	assert.NoError(t, err)

	var expected uint = 4
	assert.Equal(t, conf.PlayersInRoomCount, expected)
}

func TestLoadConfigsWithEnv(t *testing.T) {
	os.Setenv("PLAYERS_IN_ROOM_COUNT", "5")
	conf, err := LoadConfigs("./testdata/", "testconfigs", "env")
	assert.NoError(t, err)

	var expected uint = 5
	assert.Equal(t, conf.PlayersInRoomCount, expected)
}
