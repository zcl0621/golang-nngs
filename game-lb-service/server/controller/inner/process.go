package inner

import (
	"fmt"
	"game-lb/api/node"
	"game-lb/exception"
	"game-lb/request"
	"game-lb/responses"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// initInfo
// @Summary 初始化(内部)
// @Tags 对弈
// @Param request body request.InfoRequest true "对弈"
// @Router /api/v3/game-service/inner/init [post]
// @Produce json
// @Success 200 object responses.StandardResponse
func initInfo(c *gin.Context) {
	var req request.InitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("playApp.initInfo").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	e := node.InitInfo(&req)
	if e != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage(e.Error()).
			SetFunctionName("playApp.initInfo").
			SetOriginalError(e).SetErrorCode(1))
	}
	c.JSON(http.StatusOK, &responses.StandardResponse{
		Code: 0,
		Msg:  "SUCCESS",
	})
}

// callEnd
// @Summary 结束
// @Tags 内部
// @Param request body request.CallEndRequest true "对弈"
// @Router /api/v3/game-service/inner/call-end [post]
// @Produce json
// @Success 200 object responses.StandardResponse
func callEnd(c *gin.Context) {
	var req request.CallEndRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("playApp.ownership").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	e := node.CallEnd(&req)
	if e != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage(e.Error()).
			SetFunctionName("playApp.ownership").
			SetOriginalError(e).SetErrorCode(1))
	}
	c.JSON(http.StatusOK, &responses.StandardResponse{
		Code: 0,
		Msg:  "成功",
	})
}

// areaScore
// @Summary 数目
// @Tags 对弈
// @Param request body request.OwnershipRequest true "对弈"
// @Router /api/v3/game-service/inner/area-score [post]
// @Produce json
// @Success 200 object responses.InnerAreaScoreResponse
func areaScore(c *gin.Context) {
	var req request.OwnershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("inner.areaScore").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	if req.GameId == 0 {
		req.GameId = uint(time.Now().UnixNano())
	}
	res, e := node.InnerAreaScore(&req)
	if e != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage(e.Error()).
			SetFunctionName("inner.areaScore").
			SetOriginalError(e).SetErrorCode(1))
	}
	c.JSON(http.StatusOK, res)
}

// forceReload
// @Summary 强制重载棋谱
// @Tags 对弈
// @Param request body request.ForceReloadSgfRequest true "对弈"
// @Router /api/v3/game-service/inner/force-reload [post]
// @Produce json
// @Success 200 object responses.StandardResponse
func forceReload(c *gin.Context) {
	var req request.ForceReloadSgfRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("inner.forceReload").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	e := node.ForceReloadSGF(&req)
	if e != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage(e.Error()).
			SetFunctionName("inner.forceReload").
			SetOriginalError(e).SetErrorCode(1))
	}
	c.JSON(http.StatusOK, &responses.StandardResponse{
		Code: 0,
		Msg:  "SUCCESS",
	})
}
