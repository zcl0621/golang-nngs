package inner

import (
	"github.com/gin-gonic/gin"
	"higo-game-bus/exception"
)

func InitApi(api *gin.RouterGroup) {
	api.Use(exception.ExceptionStandardHandler)
	api.POST("create", create)
	api.POST("rule", rule)
	api.POST("info", info)
	api.POST("history/student", studentHistory)
	api.POST("maintain/check", checkMaintainGame)
	api.POST("score/start", startScore)
}
