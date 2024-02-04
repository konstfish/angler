package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/konstfish/angler/shared/configs"
	"github.com/konstfish/angler/shared/monitoring"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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
		ctx, span := monitoring.Tracer.Start(c.Request.Context(), "ValidateJWT", trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()

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

		c.Next()
	}
}
