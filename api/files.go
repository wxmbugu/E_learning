package api

import (
	"io"
	"log"
	"net/http"

	"github.com/E_learning/controllers"
	"github.com/E_learning/models"
	"github.com/E_learning/token"
	"github.com/E_learning/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"go.mongodb.org/mongo-driver/mongo"
)

type videoupload struct {
	Name         string `uri:"name"`
	Subsectionid string `uri:"subsectionid"`
	Sectiontitle string `uri:"sectiontitle"`
	//Content      string `form:"Content"`
}

type ProgressReader struct {
	io.Reader
	reporter func(n int)
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.reporter(n)
	return
}

func (server Server) Uploadvideo(ctx *gin.Context) {
	var req videoupload
	if err := ctx.BindUri(&req); err != nil {
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
		sess := ctx.MustGet("sess").(*session.Session)
		uploader := s3manager.NewUploader(sess)
		config, err := util.LoadConfig(".")
		if err != nil {
			log.Fatal(err)
		}
		file, header, err := ctx.Request.FormFile("file")
		if err != nil {
			log.Println(err)
		}
		head := make([]byte, 261)
		file.Read(head)
		if filetype.IsVideo(head) {
			filename := course.Author + "/" + req.Name + "/" + req.Sectiontitle + "/" + req.Subsectionid + "/" + header.Filename
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
			upload := models.Content{
				SubContent: filepath,
			}

			server.Controller.Course.UpdateSectionContent(ctx, req.Name, req.Subsectionid, req.Sectiontitle, &upload.SubContent)
			ctx.JSON(http.StatusOK, gin.H{"filepath": filepath})
		} else {
			ctx.JSON(http.StatusUnsupportedMediaType, "filetype should be video")
		}
	}

}

/*type DelContentReq struct {
	CourseName   string `uri:"name" binding:"required"`
	SubsectionId string `uri:"subsectionid" binding:"required"`
	Title        string `uri:"sectiontitle"  binding:"required"`

	type getContentRequest struct {
	Name         string `uri:"name" binding:"required"`
	SubsectionId string `uri:"subsectionid" binding:"required"`
}
}*/

func (server Server) Deletevideo(ctx *gin.Context) {
	var req getContentRequest
	content, err := server.Controller.Course.FindContent(ctx, req.Name, req.SubsectionId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong couldn't fetch data"})
		return
	}
	if content.ID.IsZero() {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}
	sess := ctx.MustGet("sess").(*session.Session)
	err = Deletevideo(sess, &server.Config.Bucketname, &content.SubContent)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, "Video deleted successfully")
}

func Deletevideo(sess *session.Session, bucket *string, item *string) error {
	svc := s3.New(sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    item,
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: bucket,
		Key:    item,
	})
	if err != nil {
		return err
	}

	return nil
}
