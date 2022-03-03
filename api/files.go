package api

import (
	"log"
	"net/http"

	"github.com/E_learning/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
)

func (server Server) Uploadvideo(ctx *gin.Context) {
	sess := ctx.MustGet("sess").(*session.Session)
	uploader := s3manager.NewUploader(sess)
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	file, header, err := ctx.Request.FormFile("video")
	if err != nil {
		log.Fatal(err)
	}
	head := make([]byte, 261)
	file.Read(head)
	if filetype.IsVideo(head) {
		filename := header.Filename
		up, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(config.Bucketname),
			ACL:    aws.String("public-read"),
			Key:    aws.String(filename),
			Body:   file,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":    "Failed to upload file",
				"uploader": up,
			})
			return
		}
		filepath := "https://" + config.Bucketname + "." + "s3-" + config.Awsregion + ".amazonaws.com/" + filename
		ctx.JSON(http.StatusOK, gin.H{"filepath": filepath})
	} else {
		ctx.JSON(http.StatusUnsupportedMediaType, "filetype should be video")
	}

}
