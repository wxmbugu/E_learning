package api

import (
	"fmt"
	"time"

	//	sess "github.com/E_learning/sessions"

	"github.com/E_learning/token"
	"github.com/E_learning/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Server struct {
	config      util.Config
	tokenMaker  token.Maker
	redisClient *redis.Client
	router      *gin.Engine
}

func NewServer(config util.Config) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetrickey)
	if err != nil {
		return nil, fmt.Errorf("couldn't Create token")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	status := redisClient.Ping()
	fmt.Println(status)
	server := Server{
		config:      config,
		tokenMaker:  tokenMaker,
		redisClient: redisClient,
	}
	server.Routes()

	return &server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) Routes() {
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.POST("/user/signup", server.CreateInstructor)
	router.POST("/user/login", server.InstructorLogin)
	authroute := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authroute.POST("/course", server.createCourse)
	authroute.DELETE("/course/delete/:id", server.deleteCourse)
	authroute.GET("/course/:id", server.findCourse)
	authroute.POST("/course/update/:id", server.updateCourse)
	authroute.GET("/courses", server.listCourses)
	authroute.POST("/course/:name", server.AddSection)
	authroute.POST("/course/:name/updatesection/:id", server.updateSection)
	authroute.DELETE("/course/:name/deletesection/:id", server.DeleteSection)
	authroute.POST("/course/:name/:sectionid", server.CreateSubSection)
	authroute.GET("/courses/:name/section/:subsectionid", server.GetSubSection)
	authroute.POST("/course/:name/update/:sectiontitle/:subsectionid", server.UpdateSubSection)
	authroute.DELETE("/course/:name/delete/:sectiontitle/:subsectionid", server.DeleteSubSection)
	server.router = router
}
