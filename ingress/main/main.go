package main

import (
	"github.com/konstfish/angler/shared/monitoring"

	"github.com/konstfish/angler/ingress/mappings"
)

func init() {
	/*configs.LoadConfig([]string{
		"MONGODB_URI",
		"REDIS_URI",
	})*/
}

func main() {
	monitoring.InitTracer("angler-ingress")

	mappings.CreateUrlMappings()
	mappings.Router.Run(":8084")
}
