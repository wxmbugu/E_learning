package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/E_learning/controllers"
	"github.com/E_learning/models"
	"github.com/E_learning/token"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseSec struct {
	Name  string `uri:"name"  binding:"required"`
	Title string `json:"title"`
}

func (server *Server) SectionWorker(context context.Context, title, id string) {

	server.wg.Add(1)
	go func() {
		defer server.wg.Done()
		ok, err := server.Controller.Course.Section(context, title, id)
		if err != nil {
			log.Println(err, title, id)
		}
		fmt.Println(ok)
	}()
}
func (server *Server) AddSection(ctx *gin.Context) {
	var req CourseSec
	//fmt.Println(req.Title)
	//var name CourseNameReq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := primitive.NewObjectID()
	section := models.Section{
		ID:    id,
		Title: req.Title,
	}

	fmt.Println(req)
	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	result, err := server.Controller.Section.AddSection(ctx, section)
	if err != nil {
		if err == controllers.ErrNoSuchDocument {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
		server.SectionWorker(ctx, req.Name, result.ID.Hex())
		server.Controller.Course.Section(ctx, req.Name, section.ID.Hex())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	server.SectionWorker(ctx, req.Name, result.ID.Hex())
	log.Println("Remove data from Redis")
	server.redisClient.Del("Courses")
	ctx.JSON(http.StatusOK, result)
}

type UpdateSectionreq struct {
	Name  string `uri:"name" binding:"required"`
	Id    string `uri:"id" binding:"required"`
	Title string `json:"Title"`
}

func (server *Server) updateSection(ctx *gin.Context) {
	//var arg models.Section
	var req UpdateSectionreq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	_, err := server.Controller.Section.FindSection(ctx, req.Id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
	} else {

		result, err := server.Controller.Section.UpdateSection(ctx, req.Id, req.Title)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if result.ModifiedCount == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
		log.Println("Remove data from Redis")
		server.redisClient.Del("Courses")
		ctx.JSON(http.StatusOK, result)
	}

}

type DelSectionreq struct {
	Name string `uri:"name" binding:"required"`
	Id   string `uri:"id" binding:"required"`
}

func (server *Server) DeleteSection(ctx *gin.Context) {
	var req DelSectionreq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	_, err := server.Controller.User.FindInstructor(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}
	section, err := server.Controller.Section.FindSection(ctx, req.Id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
		if err == controllers.ErrInvalidUser {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not Found!"})
			return
		}
	} else {
		for _, section := range section.Content {
			err := server.Controller.Content.DeleteContent(ctx, section)
			if err != nil {
				log.Println(err)
			}
		}
		err := server.Controller.Section.DeleteSection(ctx, req.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Println("Remove data from Redis")
		server.redisClient.Del("Courses")
		ctx.JSON(http.StatusOK, "Deleted successfully!")
	}
}
