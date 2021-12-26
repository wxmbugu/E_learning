package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client *mongo.Client
}

const (
	Uri    = "mongodb://localhost:27017"
	DbName = "e-learning"
)

//instantiates mongo client
func DBInstance() (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(Uri))
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
	collection := db.Client.Database(DbName).Collection(collectionName)
	return collection
}
