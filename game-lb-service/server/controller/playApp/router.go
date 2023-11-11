package playApp

import (
	"github.com/gin-gonic/gin"
)

func InitApi(api *gin.RouterGroup) {
	play := api.Group("/")
	{
		play.POST("info", info)
		play.POST("enter", enter)
		play.POST("move", move)
		play.POST("resign", resign)
		play.POST("pass", pass)
		play.POST("score", areaScore)
		play.POST("apply-score", applyAreaScore)
		play.POST("reject-score", rejectAreaScore)
		play.POST("agree-score", agreeAreaScore)
		play.POST("sgf", sgfInfo)
		play.POST("ownership", ownership)
		play.POST("apply-summation", applySummation)
		play.POST("reject-summation", rejectSummation)
		play.POST("agree-summation", agreeSummation)
		play.POST("canplay", canplay)
	}
}
