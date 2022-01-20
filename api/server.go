package api

import (
	"fmt"

	//	sess "github.com/E_learning/sessions"
	"github.com/E_learning/token"
	"github.com/E_learning/util"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Server struct {
	config       util.Config
	tokenMaker   token.Maker
	redisClient  *redis.Client
	redisSession *redisStore.Store
	router       *gin.Engine
}

func NewServer(config util.Config) (*Server, error) {
	store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
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
		config:       config,
		tokenMaker:   tokenMaker,
		redisClient:  redisClient,
		redisSession: &store,
	}
	server.Routes()

	return &server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) Routes() {
	router := gin.Default()
	router.Use(sessions.Sessions("E-learning_api", *server.redisSession))
	router.POST("/instructor/signup", server.CreateInstructor)
	router.POST("/signup", server.CreateStudent)
	router.POST("/instructor/login", server.InstructorLogin)
	router.POST("/student/login", server.StudentLogin)
	authroute := router.Group("/").Use(sessionMiddleware())

	authroute.POST("/course", server.createCourse)
	authroute.DELETE("/course/delete/:id", server.deleteCourse)
	authroute.GET("/course/:id", server.findCourse)
	authroute.POST("/course/update/:id", server.updateCourse)
	authroute.GET("/courses", server.listCourses)
	authroute.POST("/course/:name", server.AddSection)
	authroute.POST("/course/:name/updatesection/:id", server.updateSection)
	authroute.DELETE("/course/:name/deletesection/:id", server.DeleteSection)
	authroute.POST("/logout", server.logout)
	server.router = router
}
