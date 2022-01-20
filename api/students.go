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

type Studentsignupreq struct {
	FirstName string `json:"Firstname" binding:"required" `
	LastName  string `json:"Lastname" binding:"required"`
	UserName  string `json:"Username" binding:"required,alphanum"`
	Email     string `json:"Email" binding:"required"`
	Password  string `json:"Password" binding:"required,min=6"`
}
type StudentResp struct {
	Username  string    `json:"Username"`
	FirstName string    `json:"Firstname" binding:"required"`
	LastName  string    `json:"Lastname" binding:"required"`
	Email     string    `json:"Email"`
	CreatedAt time.Time `json:"Created_at"`
}

func StudentResponse(student models.Student) StudentResp {
	return StudentResp{
		Username:  student.UserName,
		FirstName: student.FirstName,
		LastName:  student.LastName,
		Email:     student.Email,
		CreatedAt: student.CreatedAt,
	}
}
func (server *Server) CreateStudent(ctx *gin.Context) {
	var req Studentsignupreq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashpassword, _ := util.HashPassword(req.Password)
	args := models.Student{
		ID:        primitive.NewObjectID(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserName:  req.UserName,
		Email:     req.Email,
		Password:  hashpassword,
		CreatedAt: time.Now(),
	}
	student, err := controllers.CreateStudent(ctx, &args)
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
	resp := StudentResponse(*student)
	ctx.JSON(http.StatusOK, resp)
}

type Studentloginreq struct {
	UserName string `json:"Username" binding:"required,alphanum"`
	Password string `json:"Password" binding:"required,min=6"`
}
type loginStudentResponse struct {
	Student StudentResp `json:"Student"`
}

func (server *Server) StudentLogin(ctx *gin.Context) {
	var req Studentloginreq

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	student, err := controllers.FindStudent(ctx, req.UserName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = util.CheckPassword(req.Password, student.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	//x := session.Get("username")
	_ = sess.SessionStart().Set(student.UserName, ctx)
	_ = loginStudentResponse{
		Student: StudentResponse(student),
	}
	ctx.JSON(http.StatusOK, "Logged In")

}
