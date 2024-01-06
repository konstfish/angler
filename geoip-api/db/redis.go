package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/konstfish/angler/geoip-api/controllers"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/time/rate"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

var geoIpCollection *mongo.Collection

func init() {
	geoIpCollection = GetCollection("angler", "geoip")
}

func NewRedisClient() *RedisClient {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URI"))
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	return &RedisClient{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

func (r *RedisClient) ListenForNewItems(queueName string) {
	log.Println("Listening for new items in queue", queueName)

	limiter := rate.NewLimiter(rate.Every(time.Minute/50), 50) // 50 requests per minute

	for {
		result, err := r.Client.BLPop(r.Ctx, 0, queueName).Result()
		if err != nil {
			panic(err)
		}

		log.Println(result)

		exists, err := CheckAddress(result[1])
		if err != nil {
			log.Fatal(err)
		}

		if !exists {
			if limiter.Allow() {
				go processAddress(result[1])
			} else {
				log.Println("Rate limit exceeded, requeueing")
				r.PushToQueue(queueName, result[1])
				time.Sleep(time.Second * 3)
			}
		}
	}
}

func processAddress(address string) {
	geoip, err := controllers.GetIpInfo(address)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"_id": geoip.Address}
	opts := options.Replace().SetUpsert(true)
	result, err := geoIpCollection.ReplaceOne(context.TODO(), filter, geoip, opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Added IP:", geoip.Address, result)
}

func (r *RedisClient) PushToQueue(queueName string, value string) {
	log.Println("Pushing to queue", queueName)
	r.Client.RPush(r.Ctx, queueName, value)
}
