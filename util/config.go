package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress   string `mapstructure:"STORAGE_SERVER_ADDRESS"`
	StorageAddress  string `mapstructure:"STORAGE_ADDRESS"`
	StoragePort     string `mapstructure:"STORAGE_PORT"`
	StorageProtocol string `mapstructure:"STORAGE_PROTOCOL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
