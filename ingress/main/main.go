package main

import (
	"github.com/konstfish/angler/ingress/configs"
	"github.com/konstfish/angler/ingress/mappings"
)

func init() {
	configs.LoadDotEnv()
}

func main() {
	mappings.CreateUrlMappings()
	mappings.Router.Run(":8084")
}
