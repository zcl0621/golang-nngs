package play

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"higo-game-node/api/ai"
	"higo-game-node/api/gameBus"
	"higo-game-node/api/ws"
	"higo-game-node/cache"
	"higo-game-node/config"
	"higo-game-node/database"
	"higo-game-node/handlerHistory"
	"higo-game-node/logger"
	"higo-game-node/model"
	"higo-game-node/notifyStruct"
	"higo-game-node/redisUtils"
	"higo-game-node/steam"
	"higo-game-node/utils"
	"higo-game-node/wq"
	"strconv"
	"sync"
	"time"
)

func notifyUserEnter(userId string, gameId uint, businessType string, ctx context.Context) {
	steam.NotifyUserEnterChan <- &steam.UserEnterNotify{
		GameId:       gameId,
		UserId:       userId,
		BusinessType: businessType,
		Ctx:          ctx,
	}
}

func startTimeDown(gameId uint, ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("startTimeDown", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	select {
	case <-ctx.Done():
		cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(gameId))
		if !ok {
			return
		}
		if cacheData.TimeDownIsStarted {
			return
		}
		if cacheData.IsEnd {
			return
		}
		//设置第0步的时间
		cacheData.LastMoveTime = time.Now().Unix()
		lockKey := fmt.Sprintf("game-node:timeDown:%d", gameId)
		locker, err := redisUtils.NewLocker([]redis.Cmdable{redisUtils.Pool}, redisUtils.Options{
			KeyPrefix:   "lock",
			LockTimeout: time.Duration(config.Conf.Rule.MoveTime) * time.Second,
			WaitTimeout: 0,
			WaitRetry:   0,
		})
		if err != nil {
			return
		}
		_, err = locker.Lock(lockKey)
		if err != nil {
			return
		}
		oneModel, _ := model.SelectOneGame(nil, &model.Game{Model: gorm.Model{ID: gameId}})
		if oneModel != nil {
			timeData := &gameTime{
				BlackTime:               oneModel.BlackTime,
				WhiteTime:               oneModel.WhiteTime,
				BlackByoYomi:            oneModel.BlackByoYomi,
				WhiteByoYomi:            oneModel.WhiteByoYomi,
				BlackByoYomiTime:        oneModel.BlackByoYomiTime,
				BlackDefaultByoYomiTime: oneModel.BlackByoYomiTime,
				WhiteByoYomiTime:        oneModel.WhiteByoYomiTime,
				WhiteDefaultByoYomiTime: oneModel.WhiteByoYomiTime,
			}
			wg := sync.WaitGroup{}
			wg.Add(3)
			go StartTimeDown(
				gameId,
				timeData,
				func() {
					wg.Done()
				},
			)
			go MoveTicker(gameId, timeData, func() {
				wg.Done()
			})
			go StartFirstMoveCountDown(gameId, func() { wg.Done() })
			wg.Wait()
			cacheData.TimeDownIsStarted = true
			db := database.GetInstance()
			tx := db.Begin()
			e := model.UpdateGames(tx, &model.Game{Model: gorm.Model{ID: gameId}}, &model.Game{IsStart: true, StartTime: time.Now().Unix(), IsShow: true})
			if e != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
			return
		}
	case <-ticker.C:
		return
	}
}

func notifyGameBegin(gameId uint, businessType string, ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("notifyGameBegin", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	select {
	case <-ctx.Done():
		lockKey := fmt.Sprintf("game-node:gameBegin:%d", gameId)
		locker, err := redisUtils.NewLocker([]redis.Cmdable{redisUtils.Pool}, redisUtils.Options{
			KeyPrefix:   "lock",
			LockTimeout: 2 * time.Hour,
			WaitTimeout: 0,
			WaitRetry:   0,
		})
		if err != nil {
			return
		}
		_, err = locker.Lock(lockKey)
		if err != nil {
			return
		}
		steam.NotifyGameBeginChan <- &steam.GameBeginNotify{
			GameId:       gameId,
			BusinessType: businessType,
			Ctx:          nil,
		}
	case <-ticker.C:
		return
	}
}

func aiMove(gameId uint, ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("aiMove", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(gameId))
	if !ok {
		return
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(gameId))
	if !ok {
		return
	}
	if cacheData.WinCaptured != 0 {
		if b.Step >= cacheData.MaxStep {
			return
		}
	}
	if cacheData.WinCaptured != 0 {
		if b.CapturesBy[wq.BLACK] >= cacheData.WinCaptured {
			return
		}
		if b.CapturesBy[wq.WHITE] >= cacheData.WinCaptured {
			return
		}
	}
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	select {
	case <-ctx.Done():
		var extra model.ExtraData
		var userId string
		var userHash string
		switch b.Player {
		case wq.BLACK:
			userId = cacheData.BlackUserId
			userHash = cacheData.BlackUserHash
			_ = json.Unmarshal([]byte(cacheData.BlackUserExtra), &extra)
		case wq.WHITE:
			userId = cacheData.WhiteUserId
			userHash = cacheData.WhiteUserHash
			_ = json.Unmarshal([]byte(cacheData.WhiteUserExtra), &extra)
		}
		var gameType string
		if cacheData.WinCaptured == 0 {
			gameType = "territory"
		} else {
			gameType = "capture"
		}
		request := ai.AiMoveRequest{
			Level:     extra.AiAgentLevel,
			StepTime:  extra.AiAgentStepTime,
			BoardSize: uint(cacheData.BoardSize),
			Type:      gameType,
		}
		if b.Step == 0 {
			request.SGF = ""
		} else {
			request.SGF = utils.ZIPSgf(b.GetSGF())
		}
		tryTimes := 3
		for {
			res, e := ai.MoveApi(&request)
			if e != nil {
				if tryTimes == 0 {
					res = "pass"
				} else if tryTimes < 0 {
					res = "resign"
				} else {
					tryTimes--
					time.Sleep(3 * time.Second)
					continue
				}
			}
			switch res {
			case "pass":
				reply := &PassResponse{}
				var err error
				err = passService(&PassRequest{
					BaseRequest: BaseRequest{
						GameId:   gameId,
						UserId:   userId,
						UserHash: userHash,
					},
					C: b.Player,
				}, cacheData, reply)
				if err != nil {
					logger.Logger("aiMove pass error", logger.ERROR, nil, fmt.Sprintf("err:%s", err.Error()))
				}
				return
			case "resign":
				reply := &EndResponse{}
				err := resignService(&ResignRequest{
					BaseRequest: BaseRequest{
						GameId:   gameId,
						UserId:   userId,
						UserHash: userHash,
					},
				}, cacheData, reply)
				if err != nil {
					logger.Logger("aiMove resign error", logger.ERROR, nil, fmt.Sprintf("err:%s", err.Error()))
				}
				return
			case "":
				reply := &PassResponse{}
				err := passService(&PassRequest{
					BaseRequest: BaseRequest{
						GameId:   gameId,
						UserId:   userId,
						UserHash: userHash,
					},
					C: b.Player,
				}, cacheData, reply)
				if err != nil {
					logger.Logger("aiMove pass error", logger.ERROR, nil, fmt.Sprintf("err:%s", err.Error()))
				}
				return
			default:
				x, y := ai.AiMoveToXY(res, cacheData.BoardSize)
				reply := &MoveResponse{}
				time.Sleep(time.Second * 3)
				var err error
				err = moveService(&MoveRequest{
					BaseRequest: BaseRequest{
						GameId:   gameId,
						UserId:   userId,
						UserHash: userHash,
					},
					X: x,
					Y: y,
					C: b.Player,
				}, cacheData, reply)
				if err == nil {
					return
				} else {
					logger.Logger("aiMove move error", logger.ERROR, nil, fmt.Sprintf("err:%s", err.Error()))
					if err.Error() != "落子失败" {
						return
					}
					break
				}
			}
		}
	case <-ticker.C:
		return
	}
}

func sendMove(gameId uint, x int, y int, c wq.Colour, turn wq.Colour, step int, pointerHash string, ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("sendMove", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	d := notifyStruct.GameMoveMsg{
		GameId:      gameId,
		X:           x,
		Y:           y,
		C:           c,
		Turn:        turn,
		Step:        step,
		PointerHash: pointerHash,
	}
	ws.GroupPublishApiChan <- &ws.GroupPublishData{
		GroupId: fmt.Sprintf("game:%d", gameId),
		Message: notifyStruct.MakeWsGameMoveMsg(&d),
		Ctx:     ctx,
	}
}

func sendPass(gameId uint, c wq.Colour, step int, pointerHash string, ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("sendPass", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	d := notifyStruct.GamePassMsg{
		GameId:      gameId,
		C:           c,
		Step:        step,
		PointerHash: pointerHash,
	}
	ws.GroupPublishApiChan <- &ws.GroupPublishData{
		GroupId: fmt.Sprintf("game:%d", gameId),
		Message: notifyStruct.MakeWsGamePassMsg(&d),
		Ctx:     ctx,
	}
}

func sendEnd(gameId uint, win int, winResult string, ctx context.Context, from string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("sendEnd", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	d := notifyStruct.GameEndMsg{
		GameId:    gameId,
		Win:       win,
		WinResult: winResult,
		From:      from,
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(gameId))
	if ok {
		d.Step = b.Step
		d.BCaptured = b.CapturesBy[wq.BLACK]
		d.WCaptured = b.CapturesBy[wq.WHITE]
		d.BScore = b.BScore
		d.WScore = b.WScore
		d.ControversyCount = b.ControversyCount
		if len(b.OwnerShip) == b.Size*b.Size {
			d.OwnerShip = *wq.GetColorOwnerShip(b.OwnerShip, b.Size)
		}
	}
	ws.GroupPublishApiChan <- &ws.GroupPublishData{
		GroupId: fmt.Sprintf("game:%d", gameId),
		Message: notifyStruct.MakeWsGameEndMsg(&d),
		Ctx:     ctx,
	}
}

func notifyGameEnd(gameId uint, win int, winResult string, cacheData *cache.CacheData, b *wq.Board, ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("notifyGameEnd", logger.ERROR, nil, fmt.Sprintf("err:%s", err))
		}
	}()
	d := &steam.GameEndNotify{
		GameId:        gameId,
		BusinessType:  cacheData.BusinessType,
		Win:           win,
		WinResult:     winResult,
		BlackUserId:   cacheData.BlackUserId,
		BlackUserType: cacheData.BlackUserType,
		WhiteUserId:   cacheData.WhiteUserId,
		WhiteUserType: cacheData.WhiteUserType,
		WinCaptured:   cacheData.WinCaptured,
		Ctx:           ctx,
	}
	if b != nil {
		d.BScore = b.BScore
		d.WScore = b.WScore
		d.BCaptured = b.CapturesBy[wq.BLACK]
		d.WCaptured = b.CapturesBy[wq.WHITE]
		d.Step = b.Step
	}
	gameModel, e := model.SelectOneGame(nil, &model.Game{Model: gorm.Model{ID: gameId}})
	if e == nil {
		d.BlackUserNickName = gameModel.BlackUserNickName
		d.BlackUserAvatar = gameModel.BlackUserAvatar
		d.WhiteUserNickName = gameModel.WhiteUserNickName
		d.WhiteUserAvatar = gameModel.WhiteUserAvatar
	}
	steam.NotifyGameEndChan <- d
}

func checkCapture(gameId uint) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("checkCapture", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(gameId))
	if !ok {
		return
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(gameId))
	if !ok {
		return
	}
	if cacheData.WinCaptured == 0 {
		return
	}
	if cacheData.MaxStep > b.Step {
		if cacheData.WinCaptured > b.CapturesBy[wq.BLACK] {
			if cacheData.WinCaptured > b.CapturesBy[wq.WHITE] {
				return
			}
		}
	}
	b.Paused = true
	shutdownCh := GetShutdownChan(gameId)
	if shutdownCh != nil {
		func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Logger("checkCapture", logger.ERROR, nil, "shutdownCh is closed")
				}
			}()
			shutdownCh <- struct{}{}
		}()
	}
	db := database.GetInstance()
	tx := db.Begin()
	var win int
	var winResult string
	if b.CapturesBy[wq.BLACK] >= cacheData.WinCaptured {
		win = 1
		winResult = "B+C" + strconv.Itoa(b.CapturesBy[wq.BLACK])
	} else if b.CapturesBy[wq.WHITE] >= cacheData.WinCaptured {
		win = 2
		winResult = "W+C" + strconv.Itoa(b.CapturesBy[wq.WHITE])
	} else if b.Step >= cacheData.MaxStep {
		// 和棋
		win = 3
		winResult = "Draw"
	}
	e := model.UpdateGames(tx,
		&model.Game{Model: gorm.Model{ID: gameId}},
		&model.Game{IsEnd: true, EndTime: time.Now().Unix(), WinResult: winResult, Win: win, SGF: b.GetSGF(), Step: b.Step, BlackCaptured: b.CapturesBy[wq.BLACK], WhiteCaptured: b.CapturesBy[wq.WHITE], IsShow: true})
	if e != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	sendEnd(gameId, win, winResult, nil, "checkCapture")
	notifyGameEnd(gameId, win, winResult, cacheData, b, nil)

	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      gameId,
		HandlerTime: time.Now().Unix(),
		Handler:     "end",
		Content:     b.GetCommitSGF(),
	}}, time.Second*5)

}

func scoreEnd(ctx context.Context, tx *gorm.DB, areaScore float64, gameId uint, b *wq.Board, cacheData *cache.CacheData) (int, string, error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("scoreEnd", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	var win int
	var winResult string
	if areaScore > b.KM {
		win = 1
		winResult = fmt.Sprintf("B+%v", areaScore-b.KM)
	} else if areaScore < b.KM {
		win = 2
		winResult = fmt.Sprintf("W+%v", b.KM-areaScore)
	} else {
		win = 3
		winResult = fmt.Sprintf("Draw")
	}
	shutdownCh := GetShutdownChan(gameId)
	if shutdownCh != nil {
		select {
		case shutdownCh <- struct{}{}:
		case <-time.After(time.Second):
			logger.Logger("scoreEnd", logger.ERROR, nil, fmt.Sprintf("gameId:%d, shutdownCh is closed", gameId))
		}
	}
	moveTickerEndCh := GetMoveTickEndChan(gameId)
	if moveTickerEndCh != nil {
		select {
		case moveTickerEndCh <- struct{}{}:
		case <-time.After(time.Second):
			logger.Logger("scoreEnd", logger.ERROR, nil, fmt.Sprintf("gameId:%d, moveTickerEndCh is closed", gameId))
		}
	}
	e := model.UpdateGames(tx,
		&model.Game{Model: gorm.Model{ID: gameId}},
		&model.Game{IsEnd: true, EndTime: time.Now().Unix(), WinResult: winResult, Win: win, SGF: b.GetSGF(), Step: b.Step,
			BlackCaptured: b.CapturesBy[wq.BLACK], WhiteCaptured: b.CapturesBy[wq.WHITE], BlackScore: b.BScore, WhiteScore: b.WScore, IsShow: true})
	if e != nil {
		return 0, "", e
	}
	sendEnd(gameId, win, winResult, ctx, "scoreEnd")
	notifyGameEnd(gameId, win, winResult, cacheData, b, ctx)
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      gameId,
		HandlerTime: time.Now().Unix(),
		Handler:     "end",
		Content:     b.GetCommitSGF(),
	}}, time.Second*5)
	return win, winResult, nil
}

func areaScore(gameId uint, b *wq.Board) (float64, float64, float64, int, *[]float64, error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("areaScore", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	if b == nil {
		return 0, 0, 0, 0, nil, errors.New("棋盘不存在")
	}
	anyData := b.GetKataGoAnalysisData(gameId)
	logger.Logger("areaScore", logger.DEBUG, nil, fmt.Sprintf("start :%d", gameId))
	e := gameBus.StartAnalysisScore(anyData)
	if e != nil {
		logger.Logger("areaScore.gameBus.StartAnalysisScore", logger.ERROR, nil, fmt.Sprintf("error gameBus.StartAnalysisScore :%d err: %s", gameId, e.Error()))
		return 0, 0, 0, 0, nil, e
	}
	for i := 0; i <= 60; i++ {
		d, e := redisUtils.Get(fmt.Sprintf(gameBus.SCOREREDISRESULTKEY, fmt.Sprintf("%d", gameId)))
		if d == nil || e != nil {
			time.Sleep(time.Millisecond * 300)
			continue
		}
		var res gameBus.AnalysisScoreResult
		e = json.Unmarshal(d, &res)
		if e != nil {
			logger.Logger("areaScore.json.Unmarshal gameBus.AnalysisScoreResult", logger.ERROR, nil, fmt.Sprintf("error gameBus.json.Unmarshal gameId:%d data:%v err: %s", gameId, d, e.Error()))
			return 0, 0, 0, 0, nil, e
		}
		var analysisRes gameBus.AnalysisScoreData
		z, e := utils.UnzipString(res.Data)
		if e != nil {
			logger.Logger("areaScore.utils.UnzipString", logger.ERROR, nil, fmt.Sprintf("error gameBus.utils.UnzipString gameId:%d data:%v err: %s", gameId, res, e.Error()))
			return 0, 0, 0, 0, nil, e
		}
		e = json.Unmarshal([]byte(z), &analysisRes)
		if e != nil {
			logger.Logger("areaScore.utils.UnzipString gameBus.AnalysisScoreData", logger.ERROR, nil, fmt.Sprintf("error gameBus.gameBus.AnalysisScoreData gameId:%d data:%s err: %s", gameId, z, e.Error()))
			return 0, 0, 0, 0, nil, e
		}
		bScore, wScore, endScore, controversyCount := gameBus.GetAnalysisScore(&analysisRes.Ownership, 0)
		return bScore, wScore, endScore, controversyCount, &analysisRes.Ownership, nil
	}
	logger.Logger("areaScore.GetResult", logger.ERROR, nil, fmt.Sprintf("超时 gameBus.StartAnalysisScore :%d", gameId))
	return 0, 0, 0, 0, nil, errors.New("数目失败")
}
