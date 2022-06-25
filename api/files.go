package api

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/E_learning/controllers"
	"github.com/E_learning/models"
	"github.com/E_learning/token"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

var channelAmqp *amqp.Channel

const uploadconcurrency = 150

type videoupload struct {
	Name         string `uri:"name"`
	Subsectionid string `uri:"subsectionid"`
	Sectiontitle string `uri:"sectiontitle"`
	//Content      string `form:"Content"`
}

type CustomeReader struct {
	fp      multipart.File
	size    int64
	read    int64
	signMap map[int64]struct{}
	mux     sync.Mutex
}

func (r *CustomeReader) Read(p []byte) (int, error) {
	return r.fp.Read(p)
}

func (r *CustomeReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := r.fp.ReadAt(p, off)
	if err != nil {
		return n, err
	}
	r.mux.Lock()
	//ignore the first signature call
	if _, ok := r.signMap[off]; ok {
		//Got lenght have read(or means has uploaded),and you can construct your message
		r.read += int64(n)
		fmt.Printf("\r[total read:%d    progress:%d%%]\n", r.read, int(float32(r.read*100)/float32(r.size)))
	} else {
		r.signMap[off] = struct{}{}
	}
	r.mux.Unlock()
	return n, err
}

func (r *CustomeReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}

func producer(filename string, rabbitmquri string) {
	amqpConnection, err := amqp.Dial(rabbitmquri)
	if err != nil {
		log.Fatal(err)
	}
	channelAmqp, _ = amqpConnection.Channel()
	q, err := channelAmqp.QueueDeclare(
		"upload", // name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(q.Name, rabbitmquri)
	//data, _ := json.Marshal()
	err = channelAmqp.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "multipart/form",
			Body:         []byte(filename),
		})
	if err != nil {
		log.Fatal(err)
	}
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
		uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
			u.PartSize = 5 * 1024 * 1024
			u.LeavePartsOnError = true
		})
		if err != nil {
			log.Fatal(err)
		}
		file, header, err := ctx.Request.FormFile("file")
		if err != nil {
			log.Println(err)
		}
		file1, ok := file.(*os.File)
		if !ok {
			log.Println(ok, file1)
		}
		filereader, _ := file1.Stat()
		head := make([]byte, 261)

		file.Read(head)
		if filetype.IsVideo(head) {
			reader := &CustomeReader{
				fp:      file,
				size:    filereader.Size(),
				signMap: map[int64]struct{}{},
			}
			filename := course.Author + "/" + req.Name + "/" + req.Sectiontitle + "/" + req.Subsectionid + "/" + header.Filename
			uploader.Concurrency = uploadconcurrency
			up, err := uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(server.Config.Bucketname),
				ACL:    aws.String("public-read"),
				Key:    aws.String(filename),
				Body:   reader,
			})

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error":    "Failed to upload file",
					"uploader": up,
				})
				return
			}
			filepath := "https://" + server.Config.Bucketname + "." + "s3-" + server.Config.Awsregion + ".amazonaws.com/" + filename
			upload := models.Content{
				SubContent: filepath,
			}
			producer(filepath, server.Config.Rabbitmquri)
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
