package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Domain struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Created float64            `bson:"created" json:"created"`
}

//type DomainSettings struct {}
