package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Session struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserAgent  string             `bson:"user_agent" json:"ua"`
	IP         string             `bson:"ip,omitempty"`
	TimeOrigin float64            `bson:"time_origin,omitempty" json:"to"`
}
