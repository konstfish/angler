package models

type Event struct {
	Location  Location `bson:"location" json:"loc"`
	Name      string   `bson:"name" json:"ev"`
	Time      float64  `bson:"time" json:"tm"`
	SessionId string   `bson:"session_id" json:"sid"`
}

type Location struct {
	Path     string `bson:"path" json:"pt"`
	Hash     string `bson:"hash" json:"hs"`
	Protocol string `bson:"protocol" json:"pro"`
}
