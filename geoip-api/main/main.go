package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/konstfish/angler/geoip-api/controllers"
	"github.com/konstfish/angler/geoip-api/db"
)

func main() {
	client, err := db.NewRabbitMQClient()
	if err != nil {
		log.Fatal(err)
	}

	go client.Listen("geoip", controllers.PushIpInfo)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Println("Received signal:", sig)

	log.Println("Shutting down")
}
