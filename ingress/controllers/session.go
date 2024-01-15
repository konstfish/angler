package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/angler/ingress/db"
	"github.com/konstfish/angler/ingress/models"
	"go.mongodb.org/mongo-driver/bson"
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
var redisClient *db.RedisClient

var sessionTTL time.Duration = time.Second * 60 * 3

func init() {
	sessionCollection = db.GetCollection("angler", "sessions")
	redisClient = db.ConnectRedis()
}

func PostSession(c *gin.Context) {
	var session models.Session

	if err := c.BindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed request"})
		return
	}

	session.ID = primitive.NewObjectID()
	session.IP = c.ClientIP()
	session.Domain = c.Param("domain")

	// c.JSON(200, session.ID)

	go writeCacheSession(session)
	go redisClient.PushToQueue("geoip", session.IP)

	result, err := writeSession(session)
	fmt.Println(err)
	fmt.Println(result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"id": session.ID})
}

func writeSession(session models.Session) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := sessionCollection.InsertOne(ctx, session)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func getSession(sessionId string) (models.Session, error) {
	var session models.Session

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectId, err := primitive.ObjectIDFromHex(sessionId)
	if err != nil {
		return session, err
	}

	err = sessionCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&session)
	if err != nil {
		return session, err
	}

	go writeCacheSession(session)
	return session, nil
}

func existsSession(sessionId string) bool {
	_, err := getCacheSession(sessionId)
	if err != nil {
		return false
	}

	return true
}

func writeCacheSession(session models.Session) {
	redisClient.Client.Set(context.Background(), session.ID.Hex(), session.SerializeSession(), sessionTTL)
}

func getCacheSession(sessionId string) (models.Session, error) {
	sessionJSON, err := redisClient.Client.Get(context.Background(), sessionId).Result()
	if err != nil {
		log.Println("cache miss")
		session, err := getSession(sessionId)

		return session, err
	}

	_, err = redisClient.Client.Expire(context.Background(), sessionId, sessionTTL).Result()

	var session models.Session
	json.Unmarshal([]byte(sessionJSON), &session)

	return session, nil
}
