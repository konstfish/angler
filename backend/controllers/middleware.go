package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/angler/backend/models"
	"github.com/konstfish/angler/shared/configs"
	"github.com/konstfish/angler/shared/monitoring"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// create new span for this middleware
		ctx := c.Request.Context()
		ctx, span := monitoring.Tracer.Start(ctx, "ValidateJWT")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		req, err := http.NewRequest("GET", configs.GetConfigVar("AUTH_URL"), nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create verification request"})
			return
		}

		// Inject the token into the request
		req.Header.Add("Authorization", authHeader)
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

		// Send the request to the auth service
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != 200 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		var user models.User
		if err := json.Unmarshal(body, &user); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		// set user in context
		c.Set("user", user)

		span.End()

		c.Next()
	}
}

func GetUserFromContext(c *gin.Context) models.User {
	user, _ := c.Get("user")
	return user.(models.User)
}
