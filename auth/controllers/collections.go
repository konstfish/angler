package controllers

import (
	"github.com/konstfish/angler/shared/db"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func InitCollections() {
	userCollection = db.GetCollection("angler", "users")
}
