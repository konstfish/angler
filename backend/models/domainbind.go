package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DomainBind struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`
	DomainID primitive.ObjectID `bson:"domain_id" json:"domain_id"`
}

//type DomainSettings struct {}
