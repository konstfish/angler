package mappings

import (
	"github.com/gin-gonic/gin"
	"github.com/konstfish/angler/backend/controllers"
	"github.com/konstfish/angler/shared/monitoring"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var Router *gin.Engine

func CreateUrlMappings() {
	Router = gin.Default()
	Router.Use(otelgin.Middleware(monitoring.ServiceName, otelgin.WithFilter(monitoring.FilterTraces)))

	Router.Use(controllers.Cors())

	v1 := Router.Group("/api/backend/v1")
	{
		v1.POST("/alive", controllers.ValidateJWT(), controllers.GetTemp)
	}
}
