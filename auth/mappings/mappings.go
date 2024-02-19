package mappings

import (
	"github.com/gin-gonic/gin"
	"github.com/konstfish/angler/auth/controllers"
	"github.com/konstfish/angler/shared/monitoring"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var Router *gin.Engine

func CreateUrlMappings() {
	// middleware
	Router = gin.Default()
	Router.Use(otelgin.Middleware(monitoring.ServiceName, otelgin.WithFilter(monitoring.FilterTraces)))

	Router.Use(controllers.Cors())

	// routes
	v1 := Router.Group("/api/auth/v1")
	{
		v1.POST("/register", controllers.PostRegister)
		v1.POST("/login", controllers.PostLogin)
		v1.GET("/verify", controllers.GetVerify)
	}
}
