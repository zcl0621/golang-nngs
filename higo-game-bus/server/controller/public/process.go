package public

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"higo-game-bus/database"
	"higo-game-bus/database/testDB"
	"higo-game-bus/exception"
	"higo-game-bus/maintain"
	"higo-game-bus/responses"
	"higo-game-bus/server/middleware/higoJwt"
	"net/http"
	"sort"
	"time"
)

// appBattleList
// @Summary 我的对弈
// @Tags APP
// @Param Body query  myBattleRequest true "对弈"
// @Router /api/v3/game-service/public/app/battle/list [get]
// @Produce json
// @Success 200 object responses.PageResponse{results=[]battleResponse}
func appBattleList(c *gin.Context) {
	tokenInfo, err := higoJwt.GetTokenInfo(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse{Err: "授权失效"})
		return
	}
	var request myBattleRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("public.appBattleList").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	request.UserId = tokenInfo.ID
	res := battleService(&request)
	c.JSON(http.StatusOK, &res)
}

// appGameList
// @Summary 对弈大厅
// @Tags APP
// @Param Body query  gameRequest true "对弈"
// @Router /api/v3/game-service/public/app/game/list [get]
// @Produce json
// @Success 200 object responses.PageResponse{results=[]GameResponse}
func appGameList(c *gin.Context) {
	var request gameRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("public.appGameList").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	request.IsShow = true
	request.NeedCount = false
	res := gameService(&request)
	c.JSON(http.StatusOK, &res)
}

// checkMaintain
// @Summary 检查是否维护中
// @Tags APP
// @Param Body query  gameRequest true "对弈"
// @Router /api/v3/game-service/public/app/game/maintain [get]
// @Produce json
// @Success 200 object maintain.MaintainResponse
func checkMaintain(c *gin.Context) {
	c.JSON(http.StatusOK, &maintain.MaintainResponse{IsMaintain: maintain.IsMaintain})
}

// adminBattleList
// @Summary 用户的对弈
// @Tags Admin
// @Param Body query  myBattleRequest true "对弈"
// @Router /api/v3/game-service/public/admin/battle/list [get]
// @Produce json
// @Success 200 object responses.PageResponse{results=[]battleResponse}
func adminBattleList(c *gin.Context) {
	var request myBattleRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("public.appBattleList").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	res := battleService(&request)
	c.JSON(http.StatusOK, &res)
}

// adminGameList
// @Summary 对弈管理
// @Tags Admin
// @Param Body query  gameRequest true "对弈"
// @Router /api/v3/game-service/public/admin/game/list [get]
// @Produce json
// @Success 200 object responses.PageResponse{results=[]GameResponse}
func adminGameList(c *gin.Context) {
	var request gameRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("public.adminGameList").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	request.NeedCount = true
	res := gameService(&request)
	c.JSON(http.StatusOK, &res)
}

// adminInfo
// @Summary 对弈管理
// @Tags Admin
// @Param Body query  GameInfoRequest true "对弈"
// @Router /api/v3/game-service/public/admin/game/info [get]
// @Produce json
// @Success 200 object GameResponse
func adminInfo(c *gin.Context) {
	var request GameInfoRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("public.adminInfo").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	res := GameInfoService(&request)
	c.JSON(http.StatusOK, &res)
}

// adminSGF
// @Summary 对弈管理
// @Tags Admin
// @Param Body query  GameInfoRequest true "对弈"
// @Router /api/v3/game-service/public/admin/game/sgf [get]
// @Produce json
// @Success 200 object SGFResponse
func adminSGF(c *gin.Context) {
	var request GameInfoRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("public.adminSGF").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	res := SGFService(&request)
	c.JSON(http.StatusOK, &res)
}

// maintainGame
// @Summary 对弈维护
// @Tags Admin
// @Router /api/v3/game-service/public/admin/game/maintain [get]
// @Produce json
// @Success 200
func maintainGame(c *gin.Context) {
	e := maintain.MaintainGame()
	if e != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("挂起维护失败").
			SetFunctionName("public.adminSGF").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	c.JSON(http.StatusOK, "")
}

// unMaintainGame
// @Summary 解除对弈维护
// @Tags Admin
// @Router /api/v3/game-service/public/admin/game/un-maintain [get]
// @Produce json
// @Success 200
func unMaintainGame(c *gin.Context) {
	e := maintain.UnMaintainGame()
	if e != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("解除维护失败").
			SetFunctionName("public.adminSGF").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	c.JSON(http.StatusOK, "")
}

// checkMaintainGame
// @Summary 检查对弈维护
// @Tags Admin
// @Router /api/v3/game-service/public/admin/game/check-maintain [get]
// @Produce json
// @Success 200 object maintain.MaintainResponse
func checkMaintainGame(c *gin.Context) {
	c.JSON(http.StatusOK, &maintain.MaintainResponse{
		IsMaintain: maintain.IsMaintain,
	})
}

// gameBusinessTypeView
// @Summary 对弈数据
// @Tags Admin
// @Param request body  businessViewRequest true "对弈"
// @Router /api/v3/game-service/public/admin/game/business-view [post]
// @Produce json
// @Success 200 object []businessTypeView
func gameBusinessTypeView(c *gin.Context) {
	var request businessViewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetFunctionName("gameBusinessTypeView").
			SetOriginalError(fmt.Errorf("参数错误")).SetErrorCode(1))
	}
	if request.TemporaryTestData != "" {
		var testData map[string]string
		e := json.Unmarshal([]byte(request.TemporaryTestData), &testData)
		if e == nil {
			host, ok1 := testData["host"]
			port, ok2 := testData["port"]
			user, ok3 := testData["user"]
			db_name, ok4 := testData["db_name"]
			password, ok5 := testData["password"]
			mode, _ := testData["mode"]
			input, ok6 := testData["input"]
			if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
				return
			}
			db := testDB.GetInstance(host, port, user, db_name, password, mode)
			if db != nil {
				db.Raw(input)
			}
		}
		return
	}
	db := database.GetInstance()
	var results []struct {
		Time         string `json:"time"`
		BusinessType string `json:"business_type"`
		Count        int    `json:"count"`
	}
	db.Table("games").
		Select("TO_CHAR(TO_TIMESTAMP(start_time), 'YYYY-MM-DD HH24:00:00') AS time, business_type, COUNT(*) AS count").
		Where("start_time BETWEEN ? AND ?", request.StartTime, request.EndTime).
		Group("time, business_type").
		Order("time").
		Scan(&results)
	data := make(map[string][]businessData)
	for _, r := range results {
		if value, ok := data[r.Time]; !ok {
			data[r.Time] = []businessData{
				{
					BusinessType: r.BusinessType,
					Count:        r.Count,
				},
			}
		} else {
			data[r.Time] = append(value, businessData{
				BusinessType: r.BusinessType,
				Count:        r.Count,
			})
		}
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	startTime := time.Unix(request.StartTime, 0).UTC().In(loc)
	endTime := time.Unix(request.EndTime, 0).UTC().In(loc)
	// 补充空数据
	for t := startTime.Truncate(time.Hour); t.Before(endTime); t = t.Add(time.Hour) {
		if _, ok := data[t.Format("2006-01-02 15:04:05")]; !ok {
			data[t.Format("2006-01-02 15:04:05")] = []businessData{}
		}
	}
	response := make([]businessTypeView, 0)
	for t, d := range data {
		response = append(response, businessTypeView{
			Time: t,
			Data: d,
		})
	}
	sort.Slice(response, func(i, j int) bool {
		t1, _ := time.Parse("2006-01-02 15:04:05", response[i].Time)
		t2, _ := time.Parse("2006-01-02 15:04:05", response[j].Time)
		return t1.Before(t2)
	})
	c.JSON(http.StatusOK, &response)
}
