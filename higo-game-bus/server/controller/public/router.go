package public

import (
	"github.com/gin-gonic/gin"
	"higo-game-bus/cacheUtils"
	"higo-game-bus/exception"
	"higo-game-bus/redisUtils"
	"higo-game-bus/server/middleware/adminJwt"
	"higo-game-bus/server/middleware/higoJwt"
	"time"
)

func InitPublicRouter(api *gin.RouterGroup) {
	appApi := api.Group("app")
	appApi.Use(higoJwt.NeedJwtAuth())
	appApi.Use(exception.ExceptionErrHandler)
	{
		appApi.GET(
			"battle/list",
			cacheUtils.CachePageWithToken(redisUtils.RedisStore, time.Minute*5, appBattleList),
		)
		appApi.GET(
			"game/list",
			cacheUtils.CachePageWithoutToken(redisUtils.RedisStore, time.Minute*5, appGameList),
		)
		appApi.GET("game/maintain", checkMaintain)
	}
	adminApi := api.Group("admin")
	adminApi.Use(adminJwt.JWTAuth())
	adminApi.Use(exception.ExceptionErrHandler)
	{
		adminApi.GET("game/list", adminGameList)
		adminApi.GET("game/info", adminInfo)
		adminApi.GET("battle/list", adminBattleList)
		adminApi.GET("game/sgf", adminSGF)
		adminApi.GET("game/maintain", maintainGame)
		adminApi.GET("game/un-maintain", unMaintainGame)
		adminApi.GET("game/check-maintain", checkMaintainGame)
		adminApi.POST("game/business-view", gameBusinessTypeView)
	}
}
