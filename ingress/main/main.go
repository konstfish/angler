package main

import (
	"github.com/konstfish/angler/shared/configs"
	"github.com/konstfish/angler/shared/db"
	"github.com/konstfish/angler/shared/monitoring"

	"github.com/konstfish/angler/ingress/controllers"
	"github.com/konstfish/angler/ingress/mappings"
)

func main() {
	configs.LoadConfig()

	monitoring.InitTracer("angler-ingress")

	db.InitMongo()
	db.InitRedis()

	controllers.InitCollections()

	mappings.CreateUrlMappings()
	mappings.Router.Run(":8084")
}
