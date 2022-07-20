package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/E_learning/controllers"
	"github.com/E_learning/models"
	"github.com/E_learning/token"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Coursereq struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `json:"name" binding:"required"`
	Author      string             `json:"author"`
	Description string             `json:"description" binding:"required"`
	CreatedAt   time.Time          `json:"created_at"`
}

func (server *Server) createCourse(ctx *gin.Context) {
	var req Coursereq
	//var x primitive.ObjectID
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	args := models.Course{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Author:      authPayload.Username,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}
	course, err := server.Controller.Course.CreateCourse(ctx, &args)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Index == 0 {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "A course with this title already exists"})
					return
				}
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Item wasn't written"})
		return
	}
	ctx.JSON(http.StatusOK, course)
	log.Println("Remove data from Redis")
	server.redisClient.Del("Courses")
}

type getCourseRequest struct {
	ID string `uri:"id"  binding:"required"`
}

func (server *Server) deleteCourse(ctx *gin.Context) {
	var req getCourseRequest
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	course, err := server.Controller.Course.FindCourse(ctx, req.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}

	}
	if course.Author != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account doesn't belong to the authenticated user"})
		return
	} else {
		for _, sectionid := range course.Section {
			section, _ := server.Controller.Section.FindSection(ctx, sectionid)
			for _, contentid := range section.Content {
				err := server.Controller.Section.DeleteSection(ctx, contentid)
				if err != nil {
					log.Println(err)
				}
			}
			err = server.Controller.Section.DeleteSection(ctx, sectionid)
			if err != nil {
				log.Println(err)
			}
		}
		err = server.Controller.Course.DeleteCourse(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong couldn't delete"})
			return
		}

		log.Println("Remove data from Redis")
		server.redisClient.Del("Courses")
		ctx.JSON(http.StatusOK, "Delete Course Successfull!")
	}

}

func (server *Server) findCourse(ctx *gin.Context) {
	var req getCourseRequest
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := server.Controller.Course.FindCourse(ctx, req.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong couldn't fetch data"})
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	log.Println(authPayload)
	if len(course.StudentsEnrolled) <= 0 {
		if authPayload.Username != course.Author {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		} else {
			ctx.JSON(http.StatusOK, course)
			return
		}
	} else {
		if authPayload.Username == course.Author {
			ctx.JSON(http.StatusOK, course)
			return
		} else {
			user, _ := server.Controller.User.FindInstructor(ctx, authPayload.Username)
			for i := 0; i < len(course.StudentsEnrolled); i++ {
				if course.StudentsEnrolled[i] == user.ID.Hex() {
					log.Println(course.StudentsEnrolled[i])
					ctx.JSON(http.StatusOK, course)
					return
				}

			}
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
}

type updateCourseRequest struct {
	ID          string `uri:"id"  binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (server *Server) updateCourse(ctx *gin.Context) {
	var req updateCourseRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	arg := controllers.UpdateCourseParams{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	instructor, err := server.Controller.User.FindInstructor(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}
	course, err := server.Controller.Course.FindCourse(ctx, arg.ID)
	if instructor.UserName != course.Author {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "account doesn't belong to the authenticated user"})
		return
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
	} else {
		results, err := server.Controller.Course.UpdateCourse(ctx, arg)
		if err != nil {
			if we, ok := err.(mongo.WriteException); ok {
				for _, e := range we.WriteErrors {
					if e.Index == 0 {
						ctx.JSON(http.StatusBadRequest, gin.H{"error": "A course with this title already exists"})
						return
					}
				}
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong!"})
			return
		}
		log.Println("Remove data from Redis")
		server.redisClient.Del("Courses")
		ctx.JSON(http.StatusOK, results)
	}
}

type listCoursesRequest struct {
	PageID   int64 `form:"page_id" binding:"required,min=0"`
	PageSize int64 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listCourses(ctx *gin.Context) {
	var req listCoursesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	instructor, err := server.Controller.User.FindInstructor(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}
	arg := controllers.ListCourseParams{
		Owner: instructor.UserName,
		Limit: req.PageSize,
		Skip:  (req.PageID - 1) * req.PageSize,
	}
	results, err := server.Controller.Course.ListCourses(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, results)

}

func (server *Server) CountCoursesbyUsers(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	totalcourses := server.Controller.Course.CountCoursesbyAuthor(ctx, authPayload.Username)
	ctx.JSON(http.StatusOK, totalcourses)
}

func (server *Server) ListAllCourses(ctx *gin.Context) {
	courses, err := server.Controller.Course.ListAllCourses(ctx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong!"})
		return
	}
	ctx.JSON(http.StatusOK, courses)

}

type Enrollreq struct {
	Coursetitle string `json:"title"`
	//Username    string `json:"username"`
}

func (server *Server) Enroll(ctx *gin.Context) {
	var req Enrollreq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.Controller.User.FindInstructor(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	course, err := server.Controller.Course.FindCoursebyName(ctx, req.Coursetitle)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if course.Author == user.UserName {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Your the course author can't enroll yourself"})
		return
	}
	result, err := server.Controller.Course.Enroll(ctx, req.Coursetitle, user.ID.Hex())
	fmt.Println("bozo")

	if err != nil {
		fmt.Println("washappening!", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) GetCoursesbyEnrollment(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := server.Controller.User.FindInstructor(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	courses, err := server.Controller.Course.EnrolledCourses(ctx, user.ID.Hex())
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Enrolled in any courses"})
		return
	}
	ctx.JSON(http.StatusOK, courses)
}
