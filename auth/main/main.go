package main

import (
	"github.com/konstfish/angler/shared/configs"
	"github.com/konstfish/angler/shared/monitoring"

	"github.com/konstfish/angler/auth/mappings"
)

func init() {
	configs.LoadConfig([]string{
		"MONGODB_URI",
		"REDIS_URI",
		"JWT_SECRET",
	})
}

func main() {
	monitoring.InitTracer("angler-auth")

	mappings.CreateUrlMappings()
	mappings.Router.Run(":8086")
}
