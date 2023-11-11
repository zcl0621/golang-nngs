package inner

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"higo-game-bus/exception"
	"higo-game-bus/maintain"
	requestModel "higo-game-bus/request"
	"higo-game-bus/responses"
	"higo-game-bus/score"
	"higo-game-bus/server/controller/public"
	"net/http"
)

// create
// @Summary 创建
// @Tags 内部接口
// @Param request body  createRequest true "对弈"
// @Router /api/v3/game-service/inner/create [post]
// @Produce json
// @Success 200 object responses.StandardResponse{data=createResponse}
func create(c *gin.Context) {
	var request createRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("inner.create").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	if maintain.IsMaintain {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("当前对弈维护中").
			SetFunctionName("inner.create").
			SetOriginalError(fmt.Errorf("当前对弈维护中")).SetErrorCode(1))
	}
	res := createService(&request)
	initGameInfo(request.CanStartTime, res.GameId)
	c.JSON(http.StatusOK, responses.StandardResponse{
		Code: 0,
		Msg:  "success",
		Data: &res,
	})
}

// rule
// @Summary 规则
// @Tags 内部接口
// @Param request body  ruleRequest true "规则"
// @Router /api/v3/game-service/inner/rule [post]
// @Produce json
// @Success 200 object responses.StandardResponse{data=ruleResponse}
func rule(c *gin.Context) {
	var request ruleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("inner.rule").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	res := ruleService(&request)
	c.JSON(http.StatusOK, responses.StandardResponse{
		Code: 0,
		Msg:  "success",
		Data: &res,
	})
}

// info
// @Summary 信息
// @Tags 内部接口
// @Param request body  public.GameInfoRequest true "对弈"
// @Router /api/v3/game-service/inner/info [post]
// @Produce json
// @Success 200 object responses.StandardResponse{data=public.GameResponse}
func info(c *gin.Context) {
	var request public.GameInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("inner.info").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	res := public.GameInfoService(&request)
	c.JSON(http.StatusOK, responses.StandardResponse{
		Code: 0,
		Msg:  "success",
		Data: &res,
	})
}

// studentHistory
// @Summary 学生历史数据
// @Tags 内部接口
// @Param request body studentHistoryRequest true "对弈"
// @Router /api/v3/game-service/inner/history/student [post]
// @Produce json
// @Success 200 object responses.StandardResponse{data=[]studentHistoryResponse}
func studentHistory(c *gin.Context) {
	var request studentHistoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("inner.studentHistory").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	res := studentHistoryService(&request)
	c.JSON(http.StatusOK, responses.StandardResponse{
		Code: 0,
		Msg:  "success",
		Data: &res,
	})
}

// checkMaintainGame
// @Summary 检查对弈维护
// @Tags 内部接口
// @Router /api/v3/game-service/inner/maintain/check [post]
// @Produce json
// @Success 200 object responses.StandardResponse{data=maintain.MaintainResponse}
func checkMaintainGame(c *gin.Context) {
	c.JSON(http.StatusOK, responses.StandardResponse{Code: 0,
		Msg: "success",
		Data: &maintain.MaintainResponse{
			IsMaintain: maintain.IsMaintain,
		}})
}

// startScore
// @Summary 开始数目
// @Tags 内部接口
// @Param request body requestModel.AnalysisScoreRequest true "对弈"
// @Router /api/v3/game-service/inner/score/start [post]
// @Produce json
// @Success 200 object responses.StandardResponse{}
func startScore(c *gin.Context) {
	var request requestModel.AnalysisScoreRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("inner.startScore").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	score.SetJob(&request)
	c.JSON(http.StatusOK, responses.StandardResponse{
		Code: 0,
		Msg:  "SUCCESS",
	})
}
