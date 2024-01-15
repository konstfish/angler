package models

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID         primitive.ObjectID `bson:"_id", json:"id"`
	UserAgent  string             `bson:"user_agent" json:"ua"`
	IP         string             `bson:"ip", json:"ip"`
	TimeOrigin float64            `bson:"time_origin" json:"to"`
	DeviceType int                `bson:"device_type" json:"dt"`
	Referrer   string             `bson:"referrer" json:"rf"`
	Domain     string             `bson:"domain" json:"dm"`
}

func (session *Session) SerializeSession() string {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return ""
	}

	return string(sessionJSON)
}
