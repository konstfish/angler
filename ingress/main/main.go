package main

import (
	"github.com/konstfish/angler/ingress/mappings"
	"github.com/konstfish/angler/ingress/monitoring"
)

func main() {
	monitoring.InitTracer("angler-ingress")

	mappings.CreateUrlMappings()
	mappings.Router.Run(":8084")
}
