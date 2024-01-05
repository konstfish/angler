package mappings

import (
	"github.com/gin-gonic/gin"
	"github.com/konstfish/angler/ingress/controllers"
)

var Router *gin.Engine

func CreateUrlMappings() {
	Router = gin.Default()

	Router.Use(controllers.Cors())

	v1 := Router.Group("/v1")
	{
		v1.POST("/session", controllers.PostSession)
		v1.POST("/event/:sessionId", controllers.PostEvent)
	}
}
