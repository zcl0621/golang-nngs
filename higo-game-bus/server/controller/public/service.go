package public

import (
	"fmt"
	"gorm.io/gorm"
	"higo-game-bus/database"
	"higo-game-bus/exception"
	"higo-game-bus/model"
	"higo-game-bus/redisUtils"
	"higo-game-bus/responses"
)

func battleService(request *myBattleRequest) *responses.PageResponse {
	db := database.GetInstance()
	bDB := db.Model(&model.Battle{}).Where(&model.Battle{UserId: request.UserId, IsShow: true})
	if request.BusinessType != "" {
		bDB = bDB.Where(&model.Battle{BusinessType: request.BusinessType})
	}
	if request.StartAtBegin != 0 {
		bDB = bDB.Where("started_at >= ?", request.StartAtBegin)
	}
	if request.StartAtEnd != 0 {
		bDB = bDB.Where("started_at <= ?", request.StartAtEnd)
	}
	if request.UserName != "" {
		bDB = bDB.Where("user_name = ?", request.UserName)
	}
	if request.UserAccount != "" {
		bDB = bDB.Where("user_account = ?", request.UserAccount)
	}
	if request.UserNickName != "" {
		bDB = bDB.Where("user_nick_name = ?", request.UserNickName)
	}
	if request.UserActualName != "" {
		bDB = bDB.Where("user_actual_name = ?", request.UserActualName)
	}
	var count int64
	bDB.Count(&count)
	var res []battleResponse
	bDB.Order("id desc").Limit(request.PageSize).Offset((request.Page - 1) * request.PageSize).Scan(&res)
	return &responses.PageResponse{
		Count:   count,
		Results: res,
	}
}

func gameService(request *gameRequest) *responses.PageResponse {
	db := database.GetInstance()
	dDB := db.Model(&model.Game{})
	if request.IsShow {
		dDB = dDB.Where(&model.Game{IsShow: request.IsShow})
	}
	if request.BusinessType != "" {
		dDB = dDB.Where(&model.Game{BusinessType: request.BusinessType})
	}
	if request.UserId != "" {
		userHash := fmt.Sprintf("%s:1", request.UserId)
		dDB = dDB.Where("black_user_hash = ? or white_user_hash = ?", userHash, userHash)
	}
	switch request.Type {
	case "territory":
		dDB = dDB.Where("win_captured = ?", 0)
	case "captured":
		dDB = dDB.Where("win_captured != ?", 0)
	}
	if request.StartAtBegin != 0 {
		dDB = dDB.Where("start_time >= ?", request.StartAtBegin)
	}
	if request.StartAtEnd != 0 {
		dDB = dDB.Where("start_time <= ?", request.StartAtEnd)
	}
	var count int64
	if request.NeedCount {
		dDB.Count(&count)
	} else {
		// 当是APP访问时 不查询count 降低数据库压力
		count = int64(request.PageSize) * 5
	}
	if count == 0 {
		return &responses.PageResponse{
			Count:   0,
			Results: nil,
		}
	}
	var res []GameResponse
	dDB.Order("id desc").Limit(request.PageSize).Offset((request.Page - 1) * request.PageSize).Scan(&res)
	return &responses.PageResponse{
		Count:   count,
		Results: res,
	}
}

func GameInfoService(request *GameInfoRequest) *GameResponse {
	db := database.GetInstance()
	var res GameResponse
	db.Model(&model.Game{}).Where("id = ?", request.GameId).Scan(&res)
	return &res
}

func SGFService(request *GameInfoRequest) *SGFResponse {
	var sgf string
	sgf = GetGameSgf(request.GameId)
	if sgf == "" {
		oneModel := model.SelectOneGame(nil, &model.Game{Model: gorm.Model{ID: request.GameId}})
		if oneModel == nil {
			panic(exception.StandardRuntimeBadError().
				SetOutPutMessage("查询错误").
				SetFunctionName("public.SGFService").
				SetOriginalError(fmt.Errorf("查询错误")).SetErrorCode(1))
		}
		if oneModel.SGF == "" {
			sgf = fmt.Sprintf("(;SZ[%d]KM[7.5])", oneModel.BoardSize)
		} else {
			sgf = oneModel.SGF
		}
	}
	return &SGFResponse{SGF: sgf}
}

var GAMESGFKEY = "game:sgf:%d"

func GetGameSgf(gameId uint) (sgf string) {
	data, err := redisUtils.Get(fmt.Sprintf(GAMESGFKEY, gameId))
	if err == nil {
		sgf = string(data)
	}
	return
}
