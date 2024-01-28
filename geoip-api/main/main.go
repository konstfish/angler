package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/konstfish/angler/geoip-api/controllers"
	"github.com/konstfish/angler/shared/db"
	"github.com/konstfish/angler/shared/monitoring"
)

func main() {
	monitoring.InitTracer("angler-geoip-api")

	redisClient := db.ConnectRedis()
	defer redisClient.Client.Close()

	go redisClient.ListenForNewItems("geoip", controllers.ProcessAddress)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Println("Received signal:", sig)

	db.DisconnectMongo()

	log.Println("Shutting down")
}
