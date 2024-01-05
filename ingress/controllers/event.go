package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/angler/ingress/db"
	"github.com/konstfish/angler/ingress/models"
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

var eventCollection *mongo.Collection

func init() {
	eventCollection = db.GetCollection("angler", "events")
}

func PostEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var event models.Event
	c.BindJSON(&event)

	event.SessionId = c.Param("sessionId")

	result, err := eventCollection.InsertOne(ctx, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(result.InsertedID)

	c.JSON(200, result)
}
