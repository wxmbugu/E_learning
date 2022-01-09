package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUri          string
	DbName         string
	Server_address string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	db_uri := os.Getenv("DB_URI")
	db_name := os.Getenv("DB_Name")
	server_address := os.Getenv("SERVER_ADDRESS")

	// now do something with s3 or whatever
	config := Config{
		DbUri:          db_uri,
		DbName:         db_name,
		Server_address: server_address,
	}
	return config, err
}
