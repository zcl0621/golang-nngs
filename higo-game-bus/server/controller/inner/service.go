package inner

import (
	"fmt"
	"gorm.io/gorm"
	"higo-game-bus/api/gameLb"
	"higo-game-bus/config"
	"higo-game-bus/database"
	"higo-game-bus/exception"
	"higo-game-bus/model"
	"higo-game-bus/pub"
	"regexp"
	"strconv"
	"time"
)

func createService(request *createRequest) *createResponse {
	tx := database.GetInstance().Begin()
	defer exception.ServiceErrorCatch(tx, nil, "inner.createService")
	var KM float64
	re := regexp.MustCompile(`KM\[(\d+\.?\d*)\]`)
	matches := re.FindStringSubmatch(request.SGF)
	if len(matches) == 2 {
		km, err := strconv.ParseFloat(matches[1], 64)
		if err == nil {
			KM = km
		}
	}
	g := &model.Game{
		BoardSize:           request.BoardSize,
		SGF:                 request.SGF,
		CanStartTime:        request.CanStartTime,
		NotEnterTime:        request.NotEnterTime,
		EnableMoveTime:      request.EnableMoveTime,
		IsStart:             request.IsStart,
		StartTime:           request.StartTime,
		IsEnd:               request.IsEnd,
		EndTime:             request.EndTime,
		WinCaptured:         request.WinCaptured,
		KM:                  KM,
		MaxStep:             request.MaxStep,
		TerritoryStep:       request.TerritoryStep,
		BlackReturnStone:    request.BlackReturnStone,
		Win:                 request.Win,
		WinResult:           request.WinResult,
		Step:                request.Step,
		BlackTime:           request.BlackTime,
		WhiteTime:           request.WhiteTime,
		BlackByoYomi:        request.BlackByoYomi,
		WhiteByoYomi:        request.WhiteByoYomi,
		BlackByoYomiTime:    request.BlackByoYomiTime,
		WhiteByoYomiTime:    request.WhiteByoYomiTime,
		BlackUserId:         request.BlackUserId,
		BlackUserName:       request.BlackUserName,
		BlackUserAccount:    request.BlackUserAccount,
		BlackUserNickName:   request.BlackUserNickName,
		BlackUserActualName: request.BlackUserActualName,
		BlackUserAvatar:     request.BlackUserAvatar,
		BlackUserLevel:      request.BlackUserLevel,
		BlackUserType:       request.BlackUserType,
		BlackUserExtra:      request.BlackUserExtra,
		BlackUserEnter:      request.BlackUserEnter,
		BlackUserHash:       fmt.Sprintf("%s:%d", request.BlackUserId, request.BlackUserType),
		WhiteUserId:         request.WhiteUserId,
		WhiteUserName:       request.WhiteUserName,
		WhiteUserAccount:    request.WhiteUserAccount,
		WhiteUserNickName:   request.WhiteUserNickName,
		WhiteUserActualName: request.WhiteUserActualName,
		WhiteUserAvatar:     request.WhiteUserAvatar,
		WhiteUserLevel:      request.WhiteUserLevel,
		WhiteUserType:       request.WhiteUserType,
		WhiteUserHash:       fmt.Sprintf("%s:%d", request.WhiteUserId, request.WhiteUserType),
		WhiteUserEnter:      request.WhiteUserEnter,
		WhiteUserExtra:      request.WhiteUserExtra,
		BusinessType:        request.BusinessType,
	}
	model.CreateGames(tx, g)
	var battleType string
	if request.WinCaptured == 0 {
		battleType = "territory"
	} else {
		battleType = "captured"
	}
	if g.BlackUserType == model.HUMANUSER {
		model.CreateBattles(tx, &model.Battle{
			GameId:                 g.ID,
			BoardSize:              request.BoardSize,
			Type:                   battleType,
			UserId:                 g.BlackUserId,
			UserName:               g.BlackUserName,
			UserAccount:            g.BlackUserAccount,
			UserNickName:           g.BlackUserNickName,
			UserActualName:         g.BlackUserActualName,
			UserAvatar:             g.BlackUserAvatar,
			UserLevel:              g.BlackUserLevel,
			UserSide:               1,
			OpponentUserId:         g.WhiteUserId,
			OpponentUserName:       g.WhiteUserName,
			OpponentUserAccount:    g.WhiteUserAccount,
			OpponentUserNickName:   g.WhiteUserNickName,
			OpponentUserActualName: g.WhiteUserActualName,
			OpponentUserAvatar:     g.WhiteUserAvatar,
			OpponentUserLevel:      g.WhiteUserLevel,
			OpponentUserSide:       2,
			BusinessType:           g.BusinessType,
		})
	}
	if g.WhiteUserType == model.HUMANUSER {
		model.CreateBattles(tx, &model.Battle{
			GameId:                 g.ID,
			BoardSize:              request.BoardSize,
			Type:                   battleType,
			UserId:                 g.WhiteUserId,
			UserName:               g.WhiteUserName,
			UserAccount:            g.WhiteUserAccount,
			UserNickName:           g.WhiteUserNickName,
			UserActualName:         g.WhiteUserActualName,
			UserAvatar:             g.WhiteUserAvatar,
			UserLevel:              g.WhiteUserLevel,
			UserSide:               2,
			OpponentUserId:         g.BlackUserId,
			OpponentUserName:       g.BlackUserName,
			OpponentUserAccount:    g.BlackUserAccount,
			OpponentUserNickName:   g.BlackUserNickName,
			OpponentUserActualName: g.BlackUserActualName,
			OpponentUserAvatar:     g.BlackUserAvatar,
			OpponentUserLevel:      g.BlackUserLevel,
			OpponentUserSide:       1,
			BusinessType:           g.BusinessType,
		})
	}
	return &createResponse{
		GameId: g.ID,
	}
}

func initGameInfo(canStartTime int64, gameId uint) {
	if canStartTime == 0 || canStartTime <= time.Now().Unix() {
		e := gameLb.InitInfoApi(&gameLb.InitRequest{GameId: gameId})
		if e != nil {
			panic(exception.StandardRuntimeBadError().
				SetOutPutMessage("创建失败").
				SetFunctionName("inner.createService").
				SetOriginalError(fmt.Errorf("创建失败")).SetErrorCode(1))
		}
	} else {
		e := pub.PubInitData(&pub.InitPubData{
			GameId:       gameId,
			CanStartTime: canStartTime,
		})
		if e != nil {
			panic(exception.StandardRuntimeBadError().
				SetOutPutMessage("创建失败").
				SetFunctionName("inner.createService").
				SetOriginalError(fmt.Errorf("创建失败")).SetErrorCode(1))
		}
	}
}

func ruleService(request *ruleRequest) *ruleResponse {
	res := ruleResponse{}
	switch request.Type {
	case 1:
		// 吃子
		switch request.BoardSize {
		case 9:
			res.SGF = "(;SZ[9]KM[0]AB[ee][ff]AW[fe][ef])"
			res.MaxStep = 50
		case 13:
			res.SGF = "(;SZ[13]KM[0]AB[fh][gg]AW[gh][fg])"
			res.MaxStep = 100
		case 19:
			res.SGF = "(;SZ[19]KM[0])"
			res.MaxStep = 100
		}
	case 2:
		//围地
		switch request.BoardSize {
		case 9:
			res.TerritoryStep = 50
			switch request.HandicapCount {
			case 0:
				res.SGF = "(;SZ[9]KM[7.5])"
			default:
				res.SGF = "(;SZ[9]KM[0])"
			}
		case 13:
			if config.RunMode == "debug" || config.RunMode == "dev" {
				res.TerritoryStep = 10
			} else {
				res.TerritoryStep = 100
			}
			switch request.HandicapCount {
			case 0:
				res.SGF = "(;SZ[13]KM[7.5])"
			case 1:
				res.SGF = "(;SZ[13]KM[0])"
			case 2:
				res.SGF = "(;SZ[13]KM[0.5]AB[dj][jd])"
				res.BlackReturnStone = 1
			case 3:
				res.SGF = "(;SZ[13]KM[0.5]AB[dj][jd][jj])"
				res.BlackReturnStone = 1.5
			case 4:
				res.SGF = "(;SZ[13]KM[0.5]AB[dj][jd][jj][dd])"
				res.BlackReturnStone = 2
			case 5:
				res.SGF = "(;SZ[13]KM[0.5]AB[jd][jj][dd][dj][gg])"
				res.BlackReturnStone = 2.5
			default:
				res.SGF = "(;SZ[13]KM[0.5]AB[jd][jj][dd][dj][gg])"
				res.BlackReturnStone = 2.5
			}
		case 19:
			if config.RunMode == "debug" || config.RunMode == "dev" {
				res.TerritoryStep = 10
			} else {
				res.TerritoryStep = 100
			}
			switch request.HandicapCount {
			case 0:
				res.SGF = "(;SZ[19]KM[7.5])"
			case 1:
				res.SGF = "(;SZ[19]KM[0])"
			case 2:
				res.SGF = "(;SZ[19]KM[0.5]AB[dp][pd])"
				res.BlackReturnStone = 1
			case 3:
				res.SGF = "(;SZ[19]KM[0.5]AB[dp][pd][dd])"
				res.BlackReturnStone = 1.5
			case 4:
				res.SGF = "(;SZ[19]KM[0.5]AB[dp][pd][dd]AB[pp])"
				res.BlackReturnStone = 2
			case 5:
				res.SGF = "(;SZ[19]KM[0.5]AB[dp][pd][dd][pp][jj])"
				res.BlackReturnStone = 2.5
			case 6:
				res.SGF = "(;SZ[19]KM[0.5]AB[dp][pd][dd][pp][dj][pj])"
				res.BlackReturnStone = 3
			case 7:
				res.SGF = "(;SZ[19]KM[0.5]AB[dp][pd][dd][pp][dj][pj][jj]))"
				res.BlackReturnStone = 3.5
			case 8:
				res.SGF = "(;SZ[19]KM[0.5]AB[dp][pd][dd][pp][dj][pj][jd][jp])"
				res.BlackReturnStone = 4
			case 9:
				res.SGF = "(;SZ[19]KM[0.5]AB[dp][pd][dd][pp][dj][pj][jd][jp][jj])"
				res.BlackReturnStone = 4.5
			default:
				res.SGF = "(;SZ[19]KM[0.5]AB[dp][pd][dd][pp][dj][pj][jd][jp][jj])"
				res.BlackReturnStone = 4.5
			}
		}
	default:
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("规则生成失败").
			SetFunctionName("inner.ruleService").
			SetOriginalError(fmt.Errorf("创建失败")).SetErrorCode(1))
	}
	return &res
}

func studentHistoryService(request *studentHistoryRequest) *[]studentHistoryResponse {
	battleModel := model.SelectBattles(nil, &model.Battle{
		UserId: fmt.Sprintf("%d", request.StudentID),
	}, func(db *gorm.DB) *gorm.DB {
		return db.Where("started_at between ? and ?", request.StartDate, request.EndDate)
	})
	if len(battleModel) == 0 {
		return nil
	}
	var res []studentHistoryResponse
	for key, v := range model.GameBusinessTypeMap {
		var d studentHistoryResponse
		d.LabelName = v
		d.Label = key
		res = append(res, d)
	}
	for i := range battleModel {
		for z := range res {
			if res[z].Label == battleModel[i].BusinessType {
				res[z].Cost += uint(battleModel[i].EndedAt - battleModel[i].StartedAt)
				if battleModel[i].UserWin == 1 {
					if battleModel[i].UserSide == 1 {
						res[z].BlackWin++
					} else {
						res[z].WhiteWin++
					}
				} else {
					if battleModel[i].UserSide == 1 {
						res[z].BlackLose++
					} else {
						res[z].WhiteLose++
					}
				}
			}
		}
	}
	return &res
}
