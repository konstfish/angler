module github.com/konstfish/angler/backend

go 1.21.1

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/konstfish/angler/shared v0.0.0-20240204164302-bcbf6795cb7a
	go.mongodb.org/mongo-driver v1.13.1
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.47.0
	golang.org/x/crypto v0.18.0
)