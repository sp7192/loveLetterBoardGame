package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configs struct {
	PlayersInRoomCount uint   `mapstructure:"PLAYERS_IN_ROOM_COUNT"`
	ClientID           uint   `mapstructure:"CLIENT_ID"`
	ClientName         string `mapstructure:"CLIENT_NAME"`

	ServerIP   string `mapstructure:"SERVER_IP"`
	ServerPort uint   `mapstructure:"SERVER_PORT"`
}

func LoadConfigs(path, name, format string) (Configs, error) {
	ret := Configs{}
	viper.SetConfigFile(fmt.Sprintf("%s/%s.%s", path, name, format))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return ret, fmt.Errorf("error loading config file: %v", err)
	}
	if err := viper.Unmarshal(&ret); err != nil {
		return ret, fmt.Errorf("error unmarshaling config: %v", err)
	}

	return ret, nil
}
