package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DbUri             string        `mapstructure:"DB_URI"`
	DbName            string        `mapstructure:"DB_Name"`
	Server_address    string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetrickey string        `mapstructure:"TOKEN_SYMETRIC_KEY"`
	Tokenduration     time.Duration `mapstructure:"TOKEN_DURATION"`
}

//const project_name = "E_learning"

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
