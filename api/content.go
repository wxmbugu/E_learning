package api

import (
	"log"
	"net/http"

	"github.com/E_learning/controllers"
	"github.com/E_learning/models"
	"github.com/E_learning/token"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func (server *Server) CreateSubSection(ctx *gin.Context) {
	var req controllers.CourseSubSection
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	course, err := controllers.FindCoursebyName(ctx, req.CourseName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
	}
	if course.Author != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": controllers.ErrInvalidUser})
		return
	} else {
		result, err := controllers.AddContent(ctx, req, authPayload.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Println("Remove data from Redis")
		server.redisClient.Del("Courses")
		ctx.JSON(http.StatusOK, result)
	}
}

type UpdateSubSectionreq struct {
	Name            string `uri:"name" binding:"required"`
	Id              string `uri:"subsectionid" binding:"required"`
	Title           string `uri:"sectiontitle"  binding:"required"`
	SubSectionTitle string `json:"Subsection_Title"`
	Content         string `json:"Content"`
}

func (server *Server) UpdateSubSection(ctx *gin.Context) {
	var req UpdateSubSectionreq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	upd := models.Content{
		SubTitle:   req.SubSectionTitle,
		SubContent: req.Content,
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	course, err := controllers.FindCoursebyName(ctx, req.Name)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if course.Author != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": controllers.ErrInvalidUser})
		return
	} else {
		content, _ := controllers.FindContent(ctx, req.Name, req.Id)
		if content.ID.IsZero() {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		result, err := controllers.UpdateSectionContent(ctx, req.Name, req.Id, req.Title, &upd)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Println("Remove data from Redis")
		server.redisClient.Del("Courses")
		ctx.JSON(http.StatusOK, result)
	}
}

type DelContentReq struct {
	CourseName   string `uri:"name" binding:"required"`
	SubsectionId string `uri:"subsectionid" binding:"required"`
	Title        string `uri:"sectiontitle"  binding:"required"`
}

func (server *Server) DeleteSubSection(ctx *gin.Context) {
	var req DelContentReq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	del := controllers.DelContent{
		CourseName:   req.CourseName,
		SubsectionId: req.SubsectionId,
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	course, err := controllers.FindCoursebyName(ctx, req.CourseName)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if course.Author != authPayload.Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": controllers.ErrInvalidUser})
		return
	} else {
		content, _ := controllers.FindContent(ctx, req.CourseName, req.SubsectionId)
		if content.ID.IsZero() {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		_, err := controllers.DeleteContent(ctx, del)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Println("Remove data from Redis")
		server.redisClient.Del("Courses")
		ctx.JSON(http.StatusOK, "Deleted successfully!")
	}
}

type getContentRequest struct {
	Name         string `uri:"name" binding:"required"`
	SubsectionId string `uri:"subsectionid" binding:"required"`
}

func (server *Server) GetSubSection(ctx *gin.Context) {
	var req getContentRequest
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	content, err := controllers.FindContent(ctx, req.Name, req.SubsectionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong couldn't fetch data"})
		return
	}
	if content.ID.IsZero() {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}
	ctx.JSON(http.StatusOK, content)
}
