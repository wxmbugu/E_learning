//Courses API
//
// This is a sample course API. The courses API is of a udemy-clone project I'm making explore.....
// Terms of service:
// there are no TOS at this moment, use at your own risk we take no responsibility
//		Schemes: http
//		Setting Up API Endpoints
//		Host: localhost:8000
//		BasePath: /
//		Version: 1.0.0
//		Contact: Stephen Wambugu <wambugumacharia35@gmail.com> https://github.com/Wambug
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
// swagger:meta
package main

import (
	"flag"
	"log"

	"github.com/E_learning/api"
	"github.com/E_learning/util"
)

// Todo list
// Loose Coupling -> Section.Content

func main() {
	env, err := util.LoadConfig(".")
	if err != nil {
		log.Print(err)
	}
	server, err := api.NewServer(env)
	address := server.Config.Server_address
	flag.StringVar(&address, "port", env.Server_address, "port address for server to run on.")
	flag.Parse()
	if err != nil {
		log.Println("Couldn't Start Server!")
	}
	err = server.Start(address)
	if err != nil {
		log.Fatal(err)
	}
}
