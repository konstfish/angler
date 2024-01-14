package db

import (
	"log"

	"github.com/konstfish/angler/geoip-api/configs"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
)

type RabbitMQClient struct {
	Connection *amqp.Connection
}

func NewRabbitMQClient() (*RabbitMQClient, error) {
	conn, err := amqp.Dial(configs.GetConfigVar("RABBITMQ_URI"))

	if err != nil {
		return nil, err
	}
	return &RabbitMQClient{Connection: conn}, nil
}

/*
limiter := rate.NewLimiter(rate.Every(time.Minute/50), 50) // 50 requests per minute
if limiter.Allow()
*/
func (c *RabbitMQClient) Listen(queueName string, handler func(msg string)) {
	_, err := rabbitmq.NewConsumer(
		c.Connection,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("consumed: %v", string(d.Body))
			go handler(string(d.Body))
			return rabbitmq.Ack
		},
		queueName,
	)
	if err != nil {
		log.Fatal(err)
	}
}

// redo this idk
func (c *RabbitMQClient) Publish(queueName string, message string) {
	publisher, err := rabbitmq.NewPublisher(c.Connection)
	if err != nil {
		log.Fatal(err)
	}

	defer publisher.Close()

	err = publisher.Publish()
}
