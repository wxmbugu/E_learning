package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/E_learning/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func UploadPdf(file, filename string) (*gridfs.UploadStream, error) {

	db, _ := db.DBInstance()

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	bucket, err := gridfs.NewBucket(db.Client.Database("e-learning"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	uploadStream, err := bucket.OpenUploadStream(
		filename,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(data)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	//log.Printf("Write file to DB was successful. File size: %d\n", fileSize)
	return uploadStream, err
}
func DownloadFile(fileName string) string {

	// For CRUD operations, here is an example
	conn, _ := db.DBInstance()
	dbe := conn.Client.Database("e-learning")
	coll := conn.OpenCollection(context.Background(), "fs.files")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var results bson.M
	err := coll.FindOne(ctx, bson.M{}).Decode(&results)
	if err != nil {
		log.Fatal(err)
	}
	// you can print out the result
	fmt.Println(results)

	bucket, _ := gridfs.NewBucket(
		dbe,
	)
	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File size to download: %v \n", dStream)
	ioutil.WriteFile(fileName, buf.Bytes(), 0600)
	return fileName
}

func DeleteFile(ctx context.Context, id string) error {

	// For CRUD operations, here is an example
	conn, _ := db.DBInstance()
	dbe := conn.Client.Database("e-learning")
	//coll := conn.OpenCollection(context.Background(), "fs.files")
	iuud, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	bucket, _ := gridfs.NewBucket(
		dbe,
	)
	err = bucket.Delete(iuud)

	return err
}

func FindFile(ctx context.Context, id string) (*mongo.Cursor, error) {

	// For CRUD operations, here is an example
	conn, _ := db.DBInstance()
	dbe := conn.Client.Database("e-learning")
	//coll := conn.OpenCollection(context.Background(), "fs.files")
	iuud, _ := primitive.ObjectIDFromHex(id)

	bucket, _ := gridfs.NewBucket(
		dbe,
	)
	cursor, err := bucket.Find(bson.M{"_id": iuud})
	if err != nil {
		log.Fatal(err)
	}
	return cursor, err
}
