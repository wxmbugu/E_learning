package main

import (
	"log"

	"github.com/E_learning/api"
	"github.com/E_learning/util"
)

func main() {
	env, err := util.LoadConfig()
	if err != nil {
		log.Print(err)
	}
	server, err := api.NewServer(env)
	if err != nil {
		log.Println("Couldn't Start Server!")
	}
	server.Start(env.Server_address)
}
