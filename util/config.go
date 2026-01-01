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

	// Bind environment variables explicitly
	viper.BindEnv("DB_DRIVER")
	viper.BindEnv("DB_SOURCE")
	viper.BindEnv("SERVER_ADDRESS")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		// If the error is that the config file was not found,
		// we ignore it and rely on environment variables.
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = nil // Clear the error since it's not critical
		} else {
			return config, err // Return other errors (syntax errors, etc.)
		}
	}

	// Unmarshal the config
	err = viper.Unmarshal(&config)
	return config, err
}
