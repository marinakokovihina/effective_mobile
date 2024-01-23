package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServerHost   string `mapstructure:"HTTP_SERVER_HOST"`
	HTTPServerPort   string `mapstructure:"HTTP_SERVER_PORT"`
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PWD"`
	PostgresDBName   string `mapstructure:"POSTGRES_DB_NAME"`
	PostgresSSLMode  string `mapstructure:"POSTGRES_SSL_MODE"`
	AgifyURL         string `mapstructure:"AGIFY_URL"`
	GenderizeURL     string `mapstructure:"GENDERIZE_URL"`
	NationalizeURL   string `mapstructure:"NATIONALIZE_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	configName, found := os.LookupEnv("CONFIG_NAME")
	if found {
		viper.SetConfigName(configName)
	} else {
		viper.SetConfigName("dev")
	}
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
