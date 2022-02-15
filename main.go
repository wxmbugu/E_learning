// Courses API
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
	"log"

	"github.com/E_learning/api"
	"github.com/E_learning/util"
)

func main() {
	env, err := util.LoadConfig(".")
	if err != nil {
		log.Print(err)
	}
	server, err := api.NewServer(env)
	if err != nil {
		log.Println("Couldn't Start Server!")
	}
	err = server.Start(env.Server_address)
	if err != nil {
		log.Fatal(err)
	}
}
