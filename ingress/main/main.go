package main

import (
	"github.com/konstfish/angler/ingress/mappings"
)

func main() {
	mappings.CreateUrlMappings()
	mappings.Router.Run(":8084")
}
