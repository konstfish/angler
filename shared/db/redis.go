package db

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/konstfish/angler/shared/configs"
	"github.com/konstfish/angler/shared/monitoring"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

type RedisQueueItem struct {
	Data    string `json:"data"`
	TraceID string `json:"traceId"`
}

func (item *RedisQueueItem) Serialize() (error, string) {
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return err, ""
	}

	return nil, string(itemJSON)
}

func (item *RedisQueueItem) Deserialize(itemJSON string) error {
	err := json.Unmarshal([]byte(itemJSON), &item)
	if err != nil {
		return err
	}

	return nil
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

func (r *RedisClient) ListenForNewItems(queueName string, handler func(ctx context.Context, msg string)) {
	log.Println("Listening for new items in queue", queueName)

	limiter := rate.NewLimiter(rate.Every(time.Minute/100), 100) // 50 requests per minute

	for {
		if limiter.Allow() {
			ctx := context.Background()
			var span trace.Span

			// pop item from queue
			result, err := r.Client.BLPop(ctx, 0, queueName).Result()
			if err != nil {
				log.Println(err)
			}

			// deserialize queue item
			var queueItem RedisQueueItem
			err = queueItem.Deserialize(result[1])
			if err != nil {
				log.Println(err)
			}

			log.Println(queueItem)

			// create span
			traceID, err := trace.TraceIDFromHex(queueItem.TraceID)
			if err != nil {
				log.Println("Invalid TraceID:", err)
			} else {
				spanContext := trace.NewSpanContext(trace.SpanContextConfig{
					TraceID: traceID,
					Remote:  true,
				})

				ctx, span = monitoring.Tracer.Start(trace.ContextWithRemoteSpanContext(ctx, spanContext), "processQueueItem")
				// span.SetAttributes(attribute.String("queue", queueName))
			}

			handler(ctx, queueItem.Data)
			if span != nil {
				span.End()
			}

		} else {
			time.Sleep(time.Second * 2)
		}
	}
}

func (r *RedisClient) PushToQueue(ctx context.Context, queueName string, value string) {
	log.Printf("Pushing %s to queue %s", value, queueName)

	log.Println(ctx)

	traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()

	// create queue item
	queueItem := RedisQueueItem{
		Data:    value,
		TraceID: traceID,
	}

	// serialize queue item
	err, item := queueItem.Serialize()
	if err != nil {
		log.Println(err)
	}

	r.Client.RPush(ctx, queueName, item)
}

func (r *RedisClient) PushToQueueWithDefaultContext(queueName string, value string) {
	r.PushToQueue(r.Ctx, queueName, value)
}
