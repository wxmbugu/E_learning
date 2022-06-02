package api

import (
	"fmt"

	//	sess "github.com/E_learning/sessions"

	"github.com/E_learning/aws"
	"github.com/E_learning/token"
	"github.com/E_learning/util"
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
	router.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
		} else {
			ctx.Next()
		}
	})
	sess := aws.ConnectAws()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("sess", sess)
		ctx.Next()
	})
	router.POST("/upload", server.Uploadvideo)
	router.POST("/user/signup", server.CreateInstructor)
	router.POST("/user/login", server.InstructorLogin)
	authroute := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authroute.GET("/total/:author", server.CountCoursesbyUsers)
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
