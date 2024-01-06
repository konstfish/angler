package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/konstfish/angler/geoip-api/configs"
	"github.com/konstfish/angler/geoip-api/db"
)

func init() {
	configs.LoadDotEnv()
}

func main() {
	redisClient := db.NewRedisClient()
	defer redisClient.Client.Close()

	go redisClient.ListenForNewItems("geoip")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Println("Received signal:", sig)

	log.Println("Shutting down")
}
