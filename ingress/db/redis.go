package db

import (
	"context"
	"log"

	"github.com/konstfish/angler/ingress/configs"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisClient() *RedisClient {
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

func (r *RedisClient) PushToQueue(queueName string, value string) {
	log.Println("Pushing to queue", queueName)
	r.Client.RPush(r.Ctx, queueName, value)
}
