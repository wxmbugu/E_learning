package controllers

import "go.mongodb.org/mongo-driver/mongo"

type Controllers struct {
	User   Instructor
	Course Course
}

func New(client *mongo.Client) Controllers {
	return Controllers{
		User: Instructor{
			client: client,
		},
		Course: Course{
			client: client,
		},
	}
}
