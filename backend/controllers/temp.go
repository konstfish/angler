package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTemp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"registered": true})
}
