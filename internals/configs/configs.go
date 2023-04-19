package configs

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Configs struct {
	PlayersInRoomCount uint `mapstructure:"PLAYERS_IN_ROOM_COUNT"`
}

func LoadConfigs(path, name, format string) (Configs, error) {
	ret := Configs{}
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(format)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("fatal error config file: default \n", err)
		os.Exit(1)
	}
	err = viper.Unmarshal(&ret)

	return ret, nil
}
