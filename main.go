package main

import (
	"log"

	"github.com/E_learning/api"
	"github.com/E_learning/util"
)

func main() {
	server := api.NewServer()
	env, err := util.LoadConfig()
	if err != nil {
		log.Print(err)
	}
	server.Start(env.Server_address)
}
