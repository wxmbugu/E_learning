package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DbUri             string        `mapstructure:"DB_URI"`
	DbName            string        `mapstructure:"DB_NAME"`
	Server_address    string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetrickey string        `mapstructure:"TOKEN_SYMETRIC_KEY"`
	Tokenduration     time.Duration `mapstructure:"TOKEN_DURATION"`
	Awsregion         string        `mapstructure:"AWS_REGION"`
	Awsaccesskey      string        `mapstructure:"AWS_ACCESS_KEY_ID"`
	Awssecretkey      string        `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	Bucketname        string        `mapstructure:"BUCKET_NAME"`
	Rabbitmquri       string        `mapstructure:"RABBITMQ_URI"`
	Rabbitmqueue      string        `mapstrucutre:"RABBITMQ_QUEUE"`
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
