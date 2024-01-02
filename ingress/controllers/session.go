package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/angler/ingress/db"
	"github.com/konstfish/angler/ingress/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
type Session struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserAgent  string             `bson:"user_agent,omitempty"`
	IP         string             `bson:"ip,omitempty"`
	TimeOrigin float64            `bson:"time_origin,omitempty"`
}

*/

var sessionCollection *mongo.Collection

func init() {
	sessionCollection = db.GetCollection("angler", "sessions")
}

func PostSession(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.Session
	c.BindJSON(&user)

	user.IP = c.ClientIP()
	user.ID = primitive.NewObjectID()

	fmt.Println(user)

	result, err := sessionCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}
