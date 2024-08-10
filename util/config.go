package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string `mapstructure:"STORAGE_SERVER_ADDRESS"`
	StorageSwaggerHost string `mapstructure:"STORAGE_SWAGGER_HOST"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
