package util

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUri             string
	DbName            string
	Server_address    string
	TokenSymmetrickey string
	Tokenduration     time.Duration
}

const project_name = "E_learning"

func LoadConfig() (Config, error) {
	re := regexp.MustCompile(`^(.*` + project_name + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	db_uri := os.Getenv("DB_URI")
	db_name := os.Getenv("DB_Name")
	server_address := os.Getenv("SERVER_ADDRESS")
	tkSymetrickey := os.Getenv("TOKEN_SYMETRIC_KEY")
	tkduration := os.Getenv("TOKEN_DURATION")
	duration, _ := time.ParseDuration(tkduration)
	// now do something with s3 or whatever
	config := Config{
		DbUri:             db_uri,
		DbName:            db_name,
		Server_address:    server_address,
		TokenSymmetrickey: tkSymetrickey,
		Tokenduration:     duration,
	}
	return config, err
}
