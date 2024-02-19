package main

import (
	"github.com/konstfish/angler/shared/configs"
	"github.com/konstfish/angler/shared/db"
	"github.com/konstfish/angler/shared/monitoring"

	"github.com/konstfish/angler/auth/controllers"
	"github.com/konstfish/angler/auth/mappings"
)

func main() {
	configs.LoadConfig("JWT_SECRET")

	monitoring.InitTracer("angler-auth")

	db.InitMongo()
	db.InitRedis()

	controllers.InitCollections()

	mappings.CreateUrlMappings()
	mappings.Router.Run(":8086")
}
