package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Business struct {
	ID     primitive.ObjectID `bson:"_id"`
	Admins []string           `json:"admins"` // array of admin user ID's
	Users  []string           `json:"users"`  // array of user ID's
}
