package main

import "github.com/E_learning/api"

func main() {
	server := api.NewServer()
	server.Start(":8000")
}
