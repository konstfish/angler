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

var eventCollection *mongo.Collection

func init() {
	eventCollection = db.GetCollection("angler", "events")
}

func PostEvent(c *gin.Context) {
	var event models.Event

	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
		return
	}

	event.SessionId = c.Param("sessionId")
	event.Time = float64(time.Now().UnixMilli())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := eventCollection.InsertOne(ctx, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(result.InsertedID)

	c.JSON(200, result)
}
