package main

import (
	"github.com/konstfish/angler/shared/configs"
	"github.com/konstfish/angler/shared/db"
	"github.com/konstfish/angler/shared/monitoring"

	"github.com/konstfish/angler/backend/controllers"
	"github.com/konstfish/angler/backend/mappings"
)

func main() {
	configs.LoadConfig("AUTH_URL")

	monitoring.InitTracer("angler-backend")

	db.InitMongo()
	db.InitRedis()

	controllers.InitCollections()

	mappings.CreateUrlMappings()
	mappings.Router.Run(":8085")
}
