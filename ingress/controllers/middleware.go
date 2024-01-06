package controllers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	}
}

func DomainReferrer() gin.HandlerFunc {
	return func(c *gin.Context) {
		referrerHeader := c.Request.Header.Get("Referer")
		target := c.Param("domain")

		referrer, err := url.Parse(referrerHeader)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid referrer URL"})
			c.Abort()
			return
		}

		log.Println(referrer.Host, target)

		if referrer.Host != target {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Referrer does not match domain"})
			c.Abort()
			return
		}

		c.Next()
	}
}
