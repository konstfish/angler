package db

import (
	"context"
	"log"
	"time"

	"github.com/konstfish/angler/shared/configs"
	"github.com/redis/go-redis/extra/redisotel/v9"
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

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.Println(err)
	}

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

func (r *RedisClient) PushToQueue(ctx context.Context, queueName string, value string) {
	log.Printf("Pushing %s to queue %s", value, queueName)
	r.Client.RPush(ctx, queueName, value)
}

func (r *RedisClient) PushToQueueWithDefaultContext(queueName string, value string) {
	r.PushToQueue(r.Ctx, queueName, value)
}