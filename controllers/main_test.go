package controllers

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/E_learning/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var controllers Controllers

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	config, err := util.LoadConfig("../.")
	if err != nil {
		log.Print(err)
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DbUri))
	if err != nil {
		log.Fatal(err)
	}
	controllers = New(client)
}
