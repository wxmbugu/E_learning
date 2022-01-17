package api

import (
	"fmt"

	"github.com/E_learning/token"
	"github.com/E_learning/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetrickey)
	if err != nil {
		return nil, fmt.Errorf("couldn't Create token")
	}
	server := Server{
		config:     config,
		tokenMaker: tokenMaker,
	}
	server.Routes()

	return &server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) Routes() {
	router := gin.Default()

	router.POST("/instructor/signup", server.CreateInstructor)
	router.POST("/signup", server.CreateStudent)
	router.POST("/instructor/login", server.InstructorLogin)
	router.POST("/student/login", server.StudentLogin)
	authroute := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authroute.POST("/course", server.createCourse)
	authroute.DELETE("/course/delete/:id", server.deleteCourse)
	authroute.GET("/course/:id", server.findCourse)
	authroute.POST("/course/update/:id", server.updateCourse)
	authroute.GET("/courses", server.listCourses)
	authroute.POST("/course/:name", server.AddSection)
	authroute.POST("/course/:name/updatesection/:id", server.updateSection)
	authroute.DELETE("/course/:name/deletesection/:id", server.DeleteSection)
	server.router = router
}
