package main

import (
	"log"

	"github.com/konstfish/angler/geoip-api/controllers"
)

func main() {
	/*redisClient := db.NewRedisClient()
	defer redisClient.Client.Close()

	go redisClient.ListenForNewItems("geoip")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Println("Received signal:", sig)

	log.Println("Shutting down")*/
	log.Println(controllers.GetIpInfo("1.1.1.1"))
}
