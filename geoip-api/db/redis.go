package db

import (
	"context"
	"log"
	"time"

	"github.com/konstfish/angler/geoip-api/configs"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func ConnectRedis() *RedisClient {
	opt, err := redis.ParseURL(configs.GetConfigVar("REDIS_URI"))
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	return &RedisClient{
		Client: rdb,
		Ctx:    context.Background(),
	}
}

func (r *RedisClient) ListenForNewItems(queueName string, handler func(msg string)) {
	log.Println("Listening for new items in queue", queueName)

	limiter := rate.NewLimiter(rate.Every(time.Minute/100), 100) // 50 requests per minute

	for {
		if limiter.Allow() {
			result, err := r.Client.BLPop(r.Ctx, 0, queueName).Result()
			if err != nil {
				panic(err)
			}
			log.Println(result)

			handler(result[1])
		} else {
			time.Sleep(time.Second * 2)
		}
	}
}

func (r *RedisClient) PushToQueue(queueName string, value string) {
	log.Println("Pushing to queue", queueName)
	r.Client.RPush(r.Ctx, queueName, value)
}
