package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/E_learning/controllers"
	"github.com/E_learning/models"
	"github.com/E_learning/token"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	//	"github.com/u2takey/go-utils/waitgroup"
	"go.mongodb.org/mongo-driver/mongo"
)

func (server *Server) ContentWorker(context context.Context, title, id string) {
	server.wg.Add(1)
	go func() {
		defer server.wg.Done()
		_, err := server.Controller.Section.Content(context, title, id)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

type Contentreq struct {
	Name         string `json:"coursetitle"`
	Content      models.Content
	Sectiontitle string
}

func (server *Server) CreateSubSection(ctx *gin.Context) {
	var req Contentreq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	course, err := server.Controller.Course.FindCoursebyName(ctx, req.Name)
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
	}
	result, err := server.Controller.Content.AddContent(ctx, req.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	server.ContentWorker(ctx, req.Sectiontitle, result.ID.Hex())
	log.Println("Remove data from Redis")
	server.redisClient.Del("Courses")
	ctx.JSON(http.StatusOK, result)

}

type UpdateSubSectionreq struct {
	Name            string `uri:"name" binding:"required"`
	Id              string `uri:"subsectionid" binding:"required"`
	Title           string `uri:"sectiontitle"  binding:"required"`
	SubSectionTitle string `json:"Subsection_Title"`
	//Content         string `json:"Content"`
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
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	course, err := server.Controller.Course.FindCoursebyName(ctx, req.Name)
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
		_, err := server.Controller.Content.FindContent(ctx, req.Id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		result, err := server.Controller.Content.UpdateContentTitle(ctx, req.Id, req.SubSectionTitle)
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
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	course, err := server.Controller.Course.FindCoursebyName(ctx, req.CourseName)
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
		content, _ := server.Controller.Content.FindContent(ctx, req.SubsectionId)
		if content.ID.IsZero() {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}

		if content.Video == "" || content.Thumbnail == "" {
			err := server.Controller.Content.DeleteContent(ctx, req.SubsectionId)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			err := server.Controller.Content.DeleteContent(ctx, req.SubsectionId)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			sess := ctx.MustGet("sess").(*session.Session)
			fmt.Println(content.Video)
			x := strings.TrimPrefix(content.Video, "https://elearning-course-videos.s3-eu-central-1.amazonaws.com/")
			fmt.Println("testing", x)
			err = Deletevideo(sess, &server.Config.Bucketname, &x)
			if err != nil {
				log.Println(err)
			}
			y := strings.TrimPrefix(content.Thumbnail, "https://elearning-course-videos.s3-eu-central-1.amazonaws.com/")
			err = Deletevideo(sess, &server.Config.Bucketname, &y)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("hjk", err)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
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
	content, err := server.Controller.Content.FindContent(ctx, req.SubsectionId)
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
