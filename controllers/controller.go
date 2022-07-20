package controllers

import "go.mongodb.org/mongo-driver/mongo"

type Controllers struct {
	User    Instructor
	Course  Course
	Section Section
	Content Content
}

func New(client *mongo.Client) Controllers {
	return Controllers{
		User: Instructor{
			client: client,
		},
		Course: Course{
			client: client,
		},
		Section: Section{
			client: client,
		},
		Content: Content{
			client: client,
		},
	}
}
