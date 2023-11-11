package playApp

import (
	"errors"
	"fmt"
	"game-lb/api/node"
	"game-lb/logger"
	"game-lb/request"
	"game-lb/responses"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func errorResponse(c *gin.Context, e error, funcName string) {
	logger.Logger(funcName, logger.ERROR, e, fmt.Sprintf("error: %s", e.Error()))
	errorCode := 1
	if e.Error() == "对弈不存在" {
		errorCode = 404
	}
	c.JSON(http.StatusBadRequest, &responses.ErrorResponse{
		Err:     e.Error(),
		ErrCode: errorCode,
	})
}

// sgfInfo
// @Summary sgf
// @Tags 对弈
// @Param request body request.InfoRequest true "对弈"
// @Router /api/v3/game-service/player/sgf [post]
// @Produce json
// @Success 200 object responses.SGFResponse
func sgfInfo(c *gin.Context) {
	var req request.InfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "sgfInfo")
		return
	}
	res, e := node.SGF(&req)
	if e != nil {
		errorResponse(c, e, "sgfInfo")
		return
	}
	c.JSON(http.StatusOK, res.Zip())
}

// info
// @Summary 基础信息
// @Tags 对弈
// @Param request body request.InfoRequest true "对弈"
// @Router /api/v3/game-service/play/info [post]
// @Produce json
// @Success 200 object responses.InfoResponse
func info(c *gin.Context) {
	var req request.InfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "info")
		return
	}
	res, e := node.Info(&req)
	if e != nil {
		errorResponse(c, e, "info")
		return
	}
	c.JSON(http.StatusOK, res.Zip())
}

// enter
// @Summary 进入
// @Tags 对弈
// @Param request body request.EnterRequest true "对弈"
// @Router /api/v3/game-service/play/enter [post]
// @Produce json
// @Success 200
func enter(c *gin.Context) {
	var req request.EnterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "enter")
		return
	}
	e := node.Enter(&req)
	if e != nil {
		errorResponse(c, e, "enter")
		return
	}
	c.JSON(http.StatusOK, "")
}

// move
// @Summary 落子
// @Tags 对弈
// @Param request body request.MoveRequest true "对弈"
// @Router /api/v3/game-service/play/move [post]
// @Produce json
// @Success 200 object responses.MoveResponse
func move(c *gin.Context) {
	var req request.MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "move")
		return
	}
	res, e := node.Move(&req)
	if e != nil {
		errorResponse(c, e, "move")
		return
	}
	c.JSON(http.StatusOK, &res)
}

// pass
// @Summary 停一手
// @Tags 对弈
// @Param request body request.PassRequest true "对弈"
// @Router /api/v3/game-service/play/pass [post]
// @Produce json
// @Success 200 object responses.PassResponse
func pass(c *gin.Context) {
	var req request.PassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "pass")
		return
	}
	res, e := node.Pass(&req)
	if e != nil {
		errorResponse(c, e, "pass")
		return
	}
	c.JSON(http.StatusOK, &res)
}

// resign
// @Summary 认输
// @Tags 对弈
// @Param request body request.ResignRequest true "对弈"
// @Router /api/v3/game-service/play/resign [post]
// @Produce json
// @Success 200 object responses.EndResponse
func resign(c *gin.Context) {
	var req request.ResignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "resign")
		return
	}
	res, e := node.Resign(&req)
	if e != nil {
		errorResponse(c, e, "resign")
		return
	}
	c.JSON(http.StatusOK, &res)
}

// areaScore
// @Summary 数目
// @Tags 对弈
// @Param request body request.InfoRequest true "对弈"
// @Router /api/v3/game-service/play/score [post]
// @Produce json
// @Success 200 object responses.EndResponse
func areaScore(c *gin.Context) {
	var req request.InfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "areaScore")
		return
	}
	res, e := node.AreaScore(&req)
	if e != nil {
		errorResponse(c, e, "areaScore")
		return
	}
	c.JSON(http.StatusOK, &res)
}

// applyAreaScore
// @Summary 申请数目
// @Tags 对弈
// @Param request body request.InfoRequest true "对弈"
// @Router /api/v3/game-service/play/apply-score [post]
// @Produce json
// @Success 200
func applyAreaScore(c *gin.Context) {
	var req request.InfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "applyAreaScore")
		return
	}
	e := node.ApplyAreaScore(&req)
	if e != nil {
		errorResponse(c, e, "applyAreaScore")
		return
	}
	c.JSON(http.StatusOK, "")
}

// agreeAreaScore
// @Summary 同意数目
// @Tags 对弈
// @Param request body request.InfoRequest true "对弈"
// @Router /api/v3/game-service/play/agree-score [post]
// @Produce json
// @Success 200 object responses.EndResponse
func agreeAreaScore(c *gin.Context) {
	var req request.InfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "agreeAreaScore")
		return
	}
	res, e := node.AgreeAreaScore(&req)
	if e != nil {
		errorResponse(c, e, "agreeAreaScore")
		return
	}
	c.JSON(http.StatusOK, &res)
}

// rejectAreaScore
// @Summary 拒绝数目
// @Tags 对弈
// @Param request body request.InfoRequest true "对弈"
// @Router /api/v3/game-service/play/reject-score [post]
// @Produce json
// @Success 200
func rejectAreaScore(c *gin.Context) {
	var req request.InfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "rejectAreaScore")
		return
	}
	e := node.RejectAreaScore(&req)
	if e != nil {
		errorResponse(c, e, "rejectAreaScore")
		return
	}
	c.JSON(http.StatusOK, "")
}

// applySummation
// @Summary 申请和棋
// @Tags 对弈
// @Param request body request.InfoRequest true "对弈"
// @Router /api/v3/game-service/play/apply-summation [post]
// @Produce json
// @Success 200
func applySummation(c *gin.Context) {
	var req request.InfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "applySummation")
		return
	}
	e := node.ApplySummation(&req)
	if e != nil {
		errorResponse(c, e, "applySummation")
		return
	}
	c.JSON(http.StatusOK, "")
}

// agreeSummation
// @Summary 同意和棋
// @Tags 对弈
// @Param request body request.InfoRequest true "对弈"
// @Router /api/v3/game-service/play/agree-summation [post]
// @Produce json
// @Success 200 object responses.EndResponse
func agreeSummation(c *gin.Context) {
	var req request.InfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "agreeSummation")
		return
	}
	res, e := node.AgreeSummation(&req)
	if e != nil {
		errorResponse(c, e, "agreeSummation")
		return
	}
	c.JSON(http.StatusOK, &res)
}

// rejectSummation
// @Summary 拒绝和棋
// @Tags 对弈
// @Param request body request.InfoRequest true "对弈"
// @Router /api/v3/game-service/play/reject-summation [post]
// @Produce json
// @Success 200
func rejectSummation(c *gin.Context) {
	var req request.InfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "rejectSummation")
		return
	}
	e := node.RejectSummation(&req)
	if e != nil {
		errorResponse(c, e, "rejectSummation")
		return
	}
	c.JSON(http.StatusOK, "")
}

// canplay
// @Summary 落子方
// @Tags 对弈
// @Param request body request.InfoRequest true "对弈"
// @Router /api/v3/game-service/play/canplay [post]
// @Produce json
// @Success 200 object responses.CanPlayResponse
func canplay(c *gin.Context) {
	var req request.InfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "canplay")
		return
	}
	res, e := node.CanPlay(&req)
	if e != nil {
		errorResponse(c, e, "canplay")
		return
	}
	c.JSON(http.StatusOK, &res)
}

// ownership
// @Summary 形势判断
// @Tags 对弈
// @Param request body request.OwnershipRequest true "对弈"
// @Router /api/v3/game-service/play/ownership [post]
// @Produce json
// @Success 200 object responses.OwnerShipResponse
func ownership(c *gin.Context) {
	var req request.OwnershipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, errors.New("参数错误:"+err.Error()), "ownership")
		return
	}
	if req.GameId == 0 {
		req.GameId = uint(time.Now().UnixNano())
	}
	res, e := node.OwnerShip(&req)
	if e != nil {
		errorResponse(c, e, "ownership")
		return
	}
	c.JSON(http.StatusOK, res)
}
