package db

import (
	"context"
	"log"
	"os"

	"github.com/konstfish/angler/geoip-api/controllers"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			geoip, err := controllers.GetIpInfo(result[1])
			if err != nil {
				log.Fatal(err)
			}

			filter := bson.M{"_id": geoip.Address}
			opts := options.Replace().SetUpsert(true)
			result, err := geoIpCollection.ReplaceOne(context.TODO(), filter, geoip, opts)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(geoip.Address, result)
		}

	}
}

/*func (r *RedisClient) PushToQueue(queueName string, value string) {
	log.Println("Pushing to queue", queueName)
	r.Client.RPush(r.Ctx, queueName, value)
}*/
