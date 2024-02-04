package main

import (
	"github.com/konstfish/angler/shared/configs"
	"github.com/konstfish/angler/shared/monitoring"

	"github.com/konstfish/angler/backend/mappings"
)

func init() {
	configs.LoadConfig([]string{
		"MONGODB_URI",
		"REDIS_URI",
		"AUTH_URL",
	})
}

func main() {
	monitoring.InitTracer("angler-backend")

	mappings.CreateUrlMappings()
	mappings.Router.Run(":8085")
}
