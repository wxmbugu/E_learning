package db

import (
	"context"
	"log"
	"time"

	"github.com/E_learning/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client *mongo.Client
}

//instantiates mongo client
func DBInstance() (*DB, error) {
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

	db := DB{
		Client: client,
	}
	return &db, nil
}

//creates collection
func (db *DB) OpenCollection(ctx context.Context, collectionName string) *mongo.Collection {
	collection := db.Client.Database("e-learning").Collection(collectionName)
	return collection
}
