package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong"})
}
