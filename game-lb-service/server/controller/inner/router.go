package inner

import (
	"game-lb/exception"
	"github.com/gin-gonic/gin"
)

func InitApi(api *gin.RouterGroup) {
	innerGroup := api.Group("/")
	{
		innerGroup.Use(exception.ExceptionStandardHandler)
		innerGroup.POST("init", initInfo)
		innerGroup.POST("call-end", callEnd)
		innerGroup.POST("area-score", areaScore)
		innerGroup.POST("force-reload", forceReload)
	}
}
