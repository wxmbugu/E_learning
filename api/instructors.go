package api

import (
	"net/http"
	"time"

	"github.com/E_learning/controllers"
	"github.com/E_learning/models"
	sess "github.com/E_learning/sessions"
	"github.com/E_learning/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Instructorsignupreq struct {
	FirstName string `json:"Firstname" binding:"required" `
	LastName  string `json:"Lastname" binding:"required"`
	UserName  string `json:"Username" binding:"required,alphanum"`
	Email     string `json:"Email" binding:"required"`
	Password  string `json:"Password" binding:"required,min=6"`
}
type InstructorResp struct {
	Username  string    `json:"Username"`
	FirstName string    `json:"Firstname" binding:"required"`
	LastName  string    `json:"Lastname" binding:"required"`
	Email     string    `json:"Email"`
	CreatedAt time.Time `json:"Created_at"`
}

func InstructorResponse(instructor models.Instructor) InstructorResp {
	return InstructorResp{
		Username:  instructor.UserName,
		FirstName: instructor.FirstName,
		LastName:  instructor.LastName,
		Email:     instructor.Email,
		CreatedAt: instructor.CreatedAt,
	}
}
func (server *Server) CreateInstructor(ctx *gin.Context) {
	var req Instructorsignupreq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashpassword, _ := util.HashPassword(req.Password)
	args := models.Instructor{
		ID:        primitive.NewObjectID(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserName:  req.UserName,
		Email:     req.Email,
		Password:  hashpassword,
		CreatedAt: time.Now(),
	}
	instructor, err := controllers.CreateInstructor(ctx, &args)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Index == 0 {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
					return
				}
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Item not written"})
		return
	}
	resp := InstructorResponse(*instructor)
	ctx.JSON(http.StatusOK, resp)
}

type Instructorloginreq struct {
	UserName string `json:"Username" binding:"required,alphanum"`
	Password string `json:"Password" binding:"required,min=6"`
}

func (server *Server) InstructorLogin(ctx *gin.Context) {
	var req Instructorloginreq

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	instructor, err := controllers.FindInstructor(ctx, req.UserName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = util.CheckPassword(req.Password, instructor.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sessf := sess.SessionStart().Set(instructor.UserName, ctx)
	sessf.Get("username", ctx)
	_ = InstructorResponse(*instructor)
	ctx.JSON(http.StatusOK, "Logged In")

}
