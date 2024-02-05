package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID    primitive.ObjectID `bson:"_id", json:"id"`
	Email string             `bson:"email" json:"email"`
}
