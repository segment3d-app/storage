package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress     string `mapstructure:"SERVER_ADDRESS"`
	ContainerName     string `mapstructure:"CONTAINER_NAME"`
	ContainerPort     string `mapstructure:"CONTAINER_PORT"`
	ContainerProtocol string `mapstructure:"CONTAINER_PROTOCOL"`
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
