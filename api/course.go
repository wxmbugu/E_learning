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
		CreatedAt:   time.Now(),
	}
	course, err := controllers.CreateCourse(ctx, &args)
	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			for _, e := range we.WriteErrors {
				// check e.Code
				if e.Index == 0 {
					ctx.JSON(http.StatusBadRequest, "A course with this title already exists")
					return
				}
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Item not written"})
		return
	}
	ctx.JSON(http.StatusOK, course)
}
