package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"higo-game-bus/config"
	_ "higo-game-bus/docs"
	"higo-game-bus/server/controller/inner"
	"higo-game-bus/server/controller/public"
	"net/http"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	innerGroup := r.Group("/api/v3/game-service/inner")
	inner.InitApi(innerGroup)
	publicGroup := r.Group("/api/v3/game-service/public")
	//rateConfig := &rateLimit.Config{
	//	Name:          "限流",
	//	Description:   "限流",
	//	Algorithm:     string(rate.TokenBucket),
	//	Interval:      time.Second,
	//	MaxRequests:   5,
	//	HeaderKeyName: "Token",
	//	ExceptionKeys: nil,
	//	Routes:        nil,
	//}
	//publicGroup.Use(rateLimit.ForGin(rateConfig))
	public.InitPublicRouter(publicGroup)
	if config.RunMode == "dev" || config.RunMode == "debug" {
		r.GET("/api/v3/game-service/public/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "")
	})
	return r
}
