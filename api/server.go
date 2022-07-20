package api

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	//	sess "github.com/E_learning/sessions"

	"github.com/E_learning/aws"
	"github.com/E_learning/controllers"
	"github.com/E_learning/token"
	"github.com/E_learning/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	Config      util.Config
	tokenMaker  token.Maker
	redisClient *redis.Client
	router      *gin.Engine
	Controller  controllers.Controllers
	wg          sync.WaitGroup
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
		Config:      config,
		tokenMaker:  tokenMaker,
		redisClient: redisClient,
		Controller:  Opendb(),
	}
	server.Routes()

	return &server, nil
}

func Opendb() controllers.Controllers {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	config, err := util.LoadConfig("../.")
	if err != nil {
		log.Print(err)
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DbUri))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	println("Connected Successfully")
	c := controllers.New(client)
	return c
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
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type,Content-Length, Accept-Encoding, Authorization")
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
	router.POST("/user/signup", server.CreateInstructor)
	router.POST("/user/login", server.InstructorLogin)

	router.GET("/courses", server.ListAllCourses)
	authroute := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authroute.GET("/course/enrolled", server.GetCoursesbyEnrollment)
	authroute.GET("/total/:author", server.CountCoursesbyUsers)
	authroute.POST("/upload/:name/:sectiontitle/:subsectionid", server.Uploadvideo)
	authroute.POST("/course", server.createCourse)
	authroute.POST("/course/:name/:sectionid", server.CreateSubSection)
	authroute.DELETE("/course/delete/:id", server.deleteCourse)
	authroute.GET("/course/:id", server.findCourse)
	authroute.POST("/course/update/:id", server.updateCourse)
	authroute.GET("/course", server.listCourses)
	authroute.POST("/enroll", server.Enroll)
	authroute.POST("/course/:name", server.AddSection)
	authroute.POST("/course/:name/updatesection/:id", server.updateSection)
	authroute.DELETE("/course/:name/deletesection/:id", server.DeleteSection)
	authroute.GET("/courses/:name/section/:subsectionid", server.GetSubSection)
	authroute.POST("/course/:name/update/:sectiontitle/:subsectionid", server.UpdateSubSection)
	authroute.DELETE("/course/:name/delete/:sectiontitle/:subsectionid", server.DeleteSubSection)
	server.router = router
}
