package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Env         string `mapstructure:"ENV"`
	Port        uint16 `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`

	RSAPrivateKey           string `mapstructure:"RSA_PRIVATE_KEY"`
	RSAPublicKey            string `mapstructure:"RSA_PUBLIC_KEY"`
	JWTTokenLifetimeInHours int    `mapstructure:"JWT_TOKEN_LIFETIME_IN_HOURS"`
	JwtInfoToken            string `mapstructure:"JWT_INFO_TOKEN"`
	MaxTimeout              int    `mapstructure:"MAX_TIMEOUT"`
}

func GetConfig() *Config {
	var config Config

	viper.AddConfigPath("../")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("unable to load config file: %v", err)
		return nil
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Println(err)
		return nil
	}

	// TODO: validate with go validator here
	return &config
}
