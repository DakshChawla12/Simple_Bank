package util

import (
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// If the error is that the config file was not found,
		// we ignore it and rely on environment variables.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return // It's a real error (e.g., syntax error)
		}
	}

	// Reset err to nil so we don't return "file not found" to TestMain
	err = viper.Unmarshal(&config)
	return config, err
}
