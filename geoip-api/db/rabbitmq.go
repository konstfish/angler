package db

import (
	"log"

	"github.com/konstfish/angler/geoip-api/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

var rabbit *amqp.Connection

func init() {
	var err error
	rabbit, err = ConnectToRabbitMQ()

	if err != nil {
		log.Fatal(err)
	}
}

func ConnectToRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial(configs.GetConfigVar("RABBITMQ_URI"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}
	return conn, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
