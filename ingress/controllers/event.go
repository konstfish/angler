package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/angler/ingress/models"
	"github.com/konstfish/angler/shared/db"
	"go.mongodb.org/mongo-driver/mongo"
)

var eventCollection *mongo.Collection

func init() {
	eventCollection = db.GetCollection("angler", "events")
}

func PostEvent(c *gin.Context) {
	ctx := c.Request.Context()

	var event models.Event

	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
		return
	}

	if !existsSession(ctx, c.Param("sessionId")) {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid session"})
		return
	}

	event.SessionId = c.Param("sessionId")
	event.Time = float64(time.Now().UnixMilli())

	result, err := eventCollection.InsertOne(ctx, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(result.InsertedID)

	c.JSON(http.StatusNoContent, nil)
}
