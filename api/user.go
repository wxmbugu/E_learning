package api

import (
	"net/http"
	"time"

	"github.com/E_learning/controllers"
	"github.com/E_learning/models"
	"github.com/E_learning/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Usersignupreq struct {
	FirstName string `json:"Firstname" binding:"required" `
	LastName  string `json:"Lastname" binding:"required"`
	UserName  string `json:"Username" binding:"required,alphanum"`
	Email     string `json:"Email" binding:"required"`
	Password  string `json:"Password" binding:"required,min=6"`
}
type UserResp struct {
	Username  string    `json:"Username"`
	FirstName string    `json:"Firstname" binding:"required"`
	LastName  string    `json:"Lastname" binding:"required"`
	Email     string    `json:"Email"`
	CreatedAt time.Time `json:"Created_at"`
}

func UserResponse(instructor models.User) UserResp {
	return UserResp{
		Username:  instructor.UserName,
		FirstName: instructor.FirstName,
		LastName:  instructor.LastName,
		Email:     instructor.Email,
		CreatedAt: instructor.CreatedAt,
	}
}
func (server *Server) CreateInstructor(ctx *gin.Context) {
	var req Usersignupreq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashpassword, _ := util.HashPassword(req.Password)
	args := models.User{
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
					ctx.JSON(http.StatusBadRequest, gin.H{"error": e.Message})
					return
				}
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Item not written"})
		return
	}
	resp := UserResponse(*instructor)
	ctx.JSON(http.StatusOK, resp)
}

type Userloginreq struct {
	UserName string `json:"Username" binding:"required"`
	Password string `json:"Password" binding:"required"`
}
type loginUserResponse struct {
	AccessToken string   `json:"access_token"`
	User        UserResp `json:"user"`
}

func (server *Server) InstructorLogin(ctx *gin.Context) {
	var req Userloginreq

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Fill in username and password"})
		return
	}
	user, err := controllers.FindInstructor(ctx, req.UserName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No Such Account!"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Login Credentials"})
		return
	}
	accessToken, err := server.tokenMaker.CreateToken(
		user.UserName,
		server.config.Tokenduration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        UserResponse(*user),
	}
	ctx.JSON(http.StatusOK, rsp)

}
