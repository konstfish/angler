package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", handlePing)

	r.Run("0.0.0.0:8084") // listen and serve on 0.0.0.0:8080
}
func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong"})
}
