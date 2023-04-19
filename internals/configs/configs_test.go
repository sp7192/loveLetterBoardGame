package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigs(t *testing.T) {
	conf, err := LoadConfigs("./testdata/", "testconfigs", "env")
	assert.NoError(t, err)

	var expected uint = 4
	assert.Equal(t, conf.PlayersInRoomCount, expected)
}
