package api

import (
	"net/http"
	"time"

	"github.com/E_learning/controllers"
	"github.com/E_learning/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (server *Server) createCourse(ctx *gin.Context) {
	var req models.Course

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	args := models.Course{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Author:      req.Author,
		Description: req.Description,
		Section:     req.Section,
		CreatedAt:   time.Now(),
	}
	course, err := controllers.CreateCourse(ctx, &args)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				if e.Index == 0 {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "A course with this title already exists"})
					return
				}
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Item not written"})
		return
	}
	ctx.JSON(http.StatusOK, course)
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
	_, err := controllers.FindCourse(ctx, req.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
	} else {
		err = controllers.DeleteCourse(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong couldn't delete"})
			return
		}
		ctx.JSON(http.StatusOK, "Delete Course Successfull!")
	}
}

func (server *Server) findCourse(ctx *gin.Context) {
	var req getCourseRequest
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	course, err := controllers.FindCourse(ctx, req.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong couldn't fetch data"})
		return
	}
	ctx.JSON(http.StatusOK, course)
}

type updateCourseRequest struct {
	ID               string           `uri:"id"  binding:"required"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	UpdateSectionReq UpdateSectionReq `json:"Section"`
}
type UpdateSectionReq struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `json:"Title"  bson:"Title,omitempty"`
	Content string             `json:"Content"  bson:"Content,omitempty"`
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
		ID:                  req.ID,
		Name:                req.Name,
		Description:         req.Description,
		UpdateSectionParams: controllers.UpdateSectionParams(req.UpdateSectionReq),
	}
	_, err := controllers.FindCourse(ctx, arg.ID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
	} else {
		results, err := controllers.UpdateCourse(ctx, arg)
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
		ctx.JSON(http.StatusOK, results)
	}
}

type listCoursesRequest struct {
	PageID   int64 `form:"page_id" binding:"required,min=0"`
	PageSize int64 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listCoursesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	arg := controllers.ListCoursesParams{
		Limit: req.PageSize,
		Skip:  (req.PageID - 1) * req.PageSize,
	}
	results, err := controllers.ListCourses(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, results)
}
