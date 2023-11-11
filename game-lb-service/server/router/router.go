package router

import (
	"game-lb/config"
	_ "game-lb/docs"
	"game-lb/server/controller/inner"
	"game-lb/server/controller/playApp"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	publicGroup := r.Group("/api/v3/game-service/play")
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
	playApp.InitApi(publicGroup)
	innerGroup := r.Group("/api/v3/game-service/inner")
	inner.InitApi(innerGroup)

	if config.RunMode == "dev" || config.RunMode == "debug" {
		r.GET("/api/v3/game-service/play/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "")
	})
	return r
}
