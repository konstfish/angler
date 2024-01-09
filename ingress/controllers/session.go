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

var sessionCollection *mongo.Collection
var redisClient *db.RedisClient

func init() {
	sessionCollection = db.GetCollection("angler", "sessions")
	redisClient = db.NewRedisClient()
}

func PostSession(c *gin.Context) {
	var session models.Session

	if err := c.BindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
		return
	}

	session.IP = c.ClientIP()
	session.Domain = c.Param("domain")

	redisClient.PushToQueue("geoip", session.IP)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := sessionCollection.InsertOne(ctx, session)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(result.InsertedID)

	c.JSON(200, result)
}
