package main

import (
	//"log"
	//"bytes"

	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	//"os"

	//"log"
	"os/exec"

	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws"

	//"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"go.mongodb.org/mongo-driver/x/mongo/driver/session"
	//vidio "github.com/AlexEidt/Vidio"
)

//"fmt"
//"log"

//ffmpeg "github.com/u2takey/ffmpeg-go"
type Config struct {
	DbUri        string `mapstructure:"DB_URI"`
	Awsregion    string `mapstructure:"AWS_REGION"`
	Awsaccesskey string `mapstructure:"AWS_ACCESS_KEY_ID"`
	Awssecretkey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	Bucketname   string `mapstructure:"BUCKET_NAME"`
	Rabbitmquri  string `mapstructure:"RABBITMQ_URI"`
	Rabbitmqueue string `mapstrucutre:"RABBITMQ_QUEUE"`
}

//const project_name = "E_learning"

var config Config

func init() {
	config, _ = LoadConfig(".")
}

func ConnectAws() *session.Session {
	Accesskeyid := config.Awsaccesskey
	Secretkeyaccess := config.Awssecretkey
	Region := config.Awsregion

	session, err := session.NewSession(&aws.Config{
		Region: &Region,
		Credentials: credentials.NewStaticCredentials(
			Accesskeyid,
			Secretkeyaccess,
			"",
		),
	})
	if err != nil {
		log.Fatal(err)
	}
	return session
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
func VideoFormatting(videourl string) []string {
	ok := strings.Split(videourl, "/")
	fmt.Println(ok)
	videofile := ok[len(ok)-1]
	cmd := exec.Command("ffmpeg",
		"-i", videourl,
		"-b:a", "128k",
		"-s", "hd1080",
		"-vcodec", "libx264",
		"-b:v", "8M",
		"-pix_fmt", "yuv420p",
		"-preset", "slow",
		"-profile:v", "baseline",
		"-movflags", "faststart",
		"-y", videofile,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
	}
	data := []string{}
	message := make(map[string]string)
	message[videourl] = videofile
	data = append(data, videofile, Thumbnail(message))
	log.Printf("[*] FINISHED FORMATING THE VIDEO IN 1920:1080 RESOLUTION %s ", videofile)
	return data
}

func Thumbnail(details map[string]string) string {
	var thumbnailname string
	for k, v := range details {
		fmt.Println(k, v)
		x := strings.Split(v, ".")
		thumbnailname = x[0]
		cmd := exec.Command("ffmpeg",
			"-y",
			"-ss", "16",
			"-i", v,
			"-frames:v", " 1",
			"-s", "720x640",
			thumbnailname+".jpg",
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + string(output))
		}
	}
	log.Printf("[*] FINISHED EXTRACTING THUMBNAIL %s :RESULT", thumbnailname)
	return thumbnailname + ".jpg"
}

var mdata map[string]string

func main() {
	//var stdout bytes.Buffer
	amqpConnection, err := amqp.Dial(config.Rabbitmquri)
	if err != nil {
		log.Fatal(err)
	}

	defer amqpConnection.Close()

	channelAmqp, err := amqpConnection.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer channelAmqp.Close()
	err = channelAmqp.Qos(
		1, //prefetch count
		0, //prefetch size
		false,
	)
	if err != nil {
		log.Fatal(err)
	}
	forever := make(chan bool)

	msgs, err := channelAmqp.Consume(
		"upload",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for d := range msgs {
			log.Printf("Receeived a video url: %s", d.Body)
			details := VideoFormatting(string(d.Body))
			fmt.Println(details)
			mdata = SplitString(string(d.Body))
			urldata, err := mapdata(mdata)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print("urldata", urldata)
			worker(details, urldata)
			fmt.Printf("user: %s video preprocessing done!\n", urldata.Author)
			d.Ack(false)
		}
	}()
	log.Printf("[*] Waiting for messages to exit precc CTRL+C")
	<-forever
}

const uploadconcurrency = 250

type CustomeReader struct {
	fp      *os.File
	size    int64
	read    int64
	signMap map[int64]struct{}
	mux     sync.Mutex
}

func (r *CustomeReader) Read(p []byte) (int, error) {
	return r.fp.Read(p)
}

///https://github.com/aws/aws-sdk-go/blob/main/example/service/s3/putObjectWithProcess/putObjWithProcess.go
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
		fmt.Printf("\r[total read:%d    progress:%d%%]", r.read, int(float32(r.read*100)/float32(r.size)))
	} else {
		r.signMap[off] = struct{}{}
	}
	r.mux.Unlock()
	return n, err
}

func (r *CustomeReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}
func Reupload(filename string, data urldata) (*string, string) {
	// The session the S3 Uploader will use
	sess := ConnectAws()
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	f, err := os.Open(filename)
	if err != nil {
		log.Panicf("failed to open file %q, %v", filename, err)
	}
	// Upload the file to S3.
	if err != nil {
		log.Fatal(err)
	}
	filereader, _ := f.Stat()
	uploader.Concurrency = uploadconcurrency
	reader := &CustomeReader{
		fp:      f,
		size:    filereader.Size(),
		signMap: map[int64]struct{}{},
	}
	newfilename := data.Author + "/" + data.Name + "/" + data.Sectiontitle + "/" + data.Subsectionid + "/" + filename
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.Bucketname),
		ACL:    aws.String("public-read"),
		Key:    aws.String(newfilename),
		Body:   reader,
	})
	if err != nil {
		log.Panicf("failed to upload file, %v", err)
	}
	log.Printf("[*] REUPLOADED THE FILE COMPLETE...%s", filename)
	return &result.Location, filename
}

func UpdateSectionContent(ctx context.Context, name, subsectionid, sectiontitle, location, filename string) (*mongo.UpdateResult, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DbUri))
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("e-learning").Collection(("Course"))
	filter := bson.D{primitive.E{Key: "Name", Value: name}}
	iuud, _ := primitive.ObjectIDFromHex(subsectionid)
	arrayFilters := options.ArrayFilters{Filters: bson.A{bson.M{"x.Title": sectiontitle}, bson.M{"y.subsectionid": iuud}}}
	upsert := true
	opts := options.UpdateOptions{
		ArrayFilters: &arrayFilters,
		Upsert:       &upsert,
	}
	subset := strings.Split(filename, ".")
	if subset[1] == "mp4" {
		update := bson.M{
			"$set": bson.M{
				//"Section.$[x].Content.$[y].Subsection_Title": arg.SubTitle,
				"Section.$[x].Content.$[y].SubContent": location,
			},
		}
		result, err := collection.UpdateOne(ctx, filter, update, &opts)
		if err != nil {
			fmt.Printf("error updating db: %+v\n", err)
		}
		log.Printf("[*] UPDATED THE SECTION CONTENT %s :RESULT", filename)
		return result, err
	} else {
		update := bson.M{
			"$set": bson.M{
				//"Section.$[x].Content.$[y].Subsection_Title": arg.SubTitle,
				"Section.$[x].Content.$[y].Thumbnail": location,
			},
		}
		result, err := collection.UpdateOne(ctx, filter, update, &opts)
		if err != nil {
			fmt.Printf("error updating db: %+v\n", err)
		}
		log.Printf("[*] UPDATED THE SECTION THUMBNAIL %s :RESULT", filename)
		return result, err
	}

}

func SplitString(s string) map[string]string {
	m := make(map[string]string)
	newstring := strings.Split(s, "/")
	for i, v := range newstring[3:] {
		switch i {
		case 0:
			m["author"] = v
		case 1:
			m["name"] = v
		case 2:
			m["sectiontitle"] = v
		case 3:
			m["subsectionid"] = v
		case 4:
			m["filename"] = s
		default:
			fmt.Println("Something went wrong")
		}
		//fmt.Println(v)
	}
	return m
}

type urldata struct {
	Author       string
	Name         string
	Sectiontitle string
	Subsectionid string
	Filename     string
}

func mapdata(m map[string]string) (urldata, error) {
	var result urldata
	//err := mapstructure.Decode(m, &result)
	jsnonbody, err := json.Marshal(m)
	if err := json.Unmarshal(jsnonbody, &result); err != nil {
		fmt.Println(err)
	}
	return result, err
}

func worker(filename []string, data urldata) (*string, string) {
	var fname string
	var location *string
	res := make(chan string)
	count := 0
	for _, files := range filename {
		count++
		go func(s string) {
			location, fname = Reupload(s, data)
			res <- fmt.Sprintf("Finished upload %s", s)

			_, err := UpdateSectionContent(
				context.Background(),
				data.Name,
				data.Subsectionid,
				data.Sectiontitle,
				*location,
				fname,
			)
			deletefiles(fname)
			if err != nil {
				log.Println(err)
			}
			//fmt.Println(ok)

		}(files)
	}
	for i := 0; i < count; i++ {
		fmt.Println(<-res)
	}
	return location, fname
}

func deletefiles(filename string) {
	cmd := exec.Command("rm", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
	}
	log.Printf("[*] DELETED %s FILE", filename)
}
