package api

import (
	"fmt"
	"net/http"

	"github.com/E_learning/controllers"
	"github.com/E_learning/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (server *Server) AddSection(ctx *gin.Context) {
	var req controllers.CourseSec
	//var name CourseNameReq
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, v := range req.Section {
		v.ID = primitive.NewObjectID()
	}
	fmt.Println(req)
	result, err := controllers.AddSection(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

type UpdateSectionreq struct {
	Name    string `uri:"name" binding:"required"`
	Id      string `uri:"id" binding:"required"`
	Title   string `json:"title"`
	Content string `json:"content"`
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
	iuud, _ := primitive.ObjectIDFromHex(req.Id)

	upd := models.Section{
		ID:      iuud,
		Title:   req.Title,
		Content: req.Content,
	}
	_, err := controllers.FindSection(ctx, req.Name, req.Id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
	} else {
		result, err := controllers.UpdateSection(ctx, req.Name, &upd)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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
	del := controllers.DelSection{
		Name: req.Name,
		Id:   req.Id,
	}
	_, err := controllers.FindSection(ctx, req.Name, req.Id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found!"})
			return
		}
	} else {
		result, err := controllers.DeleteSection(ctx, del)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
