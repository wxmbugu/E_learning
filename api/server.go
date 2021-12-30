package api

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	server := Server{}
	router := gin.Default()
	router.POST("/course", server.createCourse)
	router.DELETE("/course/delete/:id", server.deleteCourse)
	router.GET("/course/:id", server.findCourse)
	router.POST("/course/update/:id", server.updateCourse)
	router.GET("/courses", server.listAccounts)

	server.router = router
	return &server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) Routes() {
	router := gin.Default()
	router.POST("/courses", server.createCourse)
}
