package controllers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/angler/shared/monitoring"
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
		ctx := c.Request.Context()
		ctx, span := monitoring.Tracer.Start(ctx, "DomainReferrerMiddleware")

		referrerHeader := c.Request.Header.Get("Referer")
		target := c.Param("domain")

		referrer, err := url.Parse(referrerHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid referrer URL"})
			return
		}

		log.Println(referrer.Host, target)

		if referrer.Host != target {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Referrer does not match domain"})
			return
		}

		// check if the referrer is in the database
		valid, err := ValidateDomain(ctx, target)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate domain"})
			return
		}
		if !valid {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
			return
		}

		span.End()
		c.Next()
	}
}
