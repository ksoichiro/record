package config

import (
	"github.com/spf13/viper"
	"log"
)

var config *viper.Viper

// Init initializes Viper config object with yaml files.
func Init(env string) {
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	if err := config.ReadInConfig(); err != nil {
		log.Fatal("failed to read config")
	}
}

// GetConfig returns the viper config.
func GetConfig() *viper.Viper {
	return config
}
