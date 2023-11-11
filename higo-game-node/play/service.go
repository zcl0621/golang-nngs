package play

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"higo-game-node/api/ws"
	"higo-game-node/cache"
	"higo-game-node/config"
	"higo-game-node/database"
	"higo-game-node/handlerHistory"
	"higo-game-node/logger"
	"higo-game-node/model"
	"higo-game-node/notifyStruct"
	"higo-game-node/utils"
	"higo-game-node/wq"
	"time"
)

func sgfService(request *InfoRequest, res *SGFResponse) error {
	d, ok := cache.CachaDataMap.Get(cache.GetKey(request.GameId))
	if ok {
		if d.IsMove != 0 {
			count := 50
			ch := make(chan struct{}, 1)
			defer close(ch)
			go func() {
				for {
					if d.IsMove == 0 || count == 0 {
						ch <- struct{}{}
						break
					}
					count--
					time.Sleep(time.Millisecond * 80)
				}
			}()
			<-ch
		}
	}
	gSGF := cache.GetGameSgf(request.GameId)
	if gSGF == "" {
		oneModel, _ := model.SelectOneGame(nil, &model.Game{Model: gorm.Model{ID: request.GameId}})
		if oneModel == nil {
			return errors.New("对弈不存在")
		}
		if oneModel.SGF == "" {
			gSGF = fmt.Sprintf("(;SZ[%d]KM[7.5])", oneModel.BoardSize)
		} else {
			gSGF = oneModel.SGF
		}
	}
	*res = SGFResponse{
		SGF: gSGF,
	}
	return nil
}

func initService(request *InitRequest) error {
	cData := cache.MakeCacheData(request.GameId)
	_, e := wq.MakeWq(request.GameId, cData.NotEnterTime)
	if e != nil {
		return e
	}
	go WithNotEnterTicker(request.GameId)
	AddPlayingGame(request.GameId)
	return nil
}

func infoService(request *InfoRequest, res *InfoResponse) error {
	oneModel, _ := model.SelectOneGame(nil, &model.Game{Model: gorm.Model{ID: request.GameId}})
	if oneModel == nil {
		return errors.New("对弈不存在")
	}
	*res = InfoResponse{
		BoardSize:           oneModel.BoardSize,
		NotEnterTime:        oneModel.NotEnterTime,
		IsStart:             oneModel.IsStart,
		StartTime:           oneModel.StartTime,
		IsEnd:               oneModel.IsEnd,
		EndTime:             oneModel.EndTime,
		Win:                 oneModel.Win,
		WinResult:           oneModel.WinResult,
		WinCapture:          oneModel.WinCaptured,
		KM:                  oneModel.KM,
		MoveTime:            config.Conf.Rule.MoveTime,
		EnableMoveTime:      oneModel.EnableMoveTime,
		SummationCount:      config.Conf.Rule.SummationCount,
		MaxStep:             oneModel.MaxStep,
		TerritoryStep:       oneModel.TerritoryStep,
		BlackReturnStone:    oneModel.BlackReturnStone,
		Step:                oneModel.Step,
		BlackTime:           oneModel.BlackTime,
		WhiteTime:           oneModel.WhiteTime,
		BlackByoYomi:        oneModel.BlackByoYomi,
		WhiteByoYomi:        oneModel.WhiteByoYomi,
		BlackByoYomiTime:    oneModel.BlackByoYomiTime,
		WhiteByoYomiTime:    oneModel.WhiteByoYomiTime,
		BlackUserId:         oneModel.BlackUserId,
		BlackUserName:       oneModel.BlackUserName,
		BlackUserAccount:    oneModel.BlackUserAccount,
		BlackUserNickName:   oneModel.BlackUserNickName,
		BlackUserActualName: oneModel.BlackUserActualName,
		BlackUserAvatar:     oneModel.BlackUserAvatar,
		BlackUserLevel:      oneModel.BlackUserLevel,
		BlackUserHash:       oneModel.BlackUserHash,
		BlackUserType:       oneModel.BlackUserType,
		BlackUserEnter:      false,
		BlackUserOnline:     false,
		WhiteUserId:         oneModel.WhiteUserId,
		WhiteUserName:       oneModel.WhiteUserName,
		WhiteUserAccount:    oneModel.WhiteUserAccount,
		WhiteUserNickName:   oneModel.WhiteUserNickName,
		WhiteUserActualName: oneModel.WhiteUserActualName,
		WhiteUserAvatar:     oneModel.WhiteUserAvatar,
		WhiteUserLevel:      oneModel.WhiteUserLevel,
		WhiteUserHash:       oneModel.WhiteUserHash,
		WhiteUserType:       oneModel.WhiteUserType,
		WhiteUserEnter:      false,
		WhiteUserOnline:     false,
		BusinessType:        oneModel.BusinessType,
		BlackCaptured:       oneModel.BlackCaptured,
		WhiteCaptured:       oneModel.WhiteCaptured,
		BlackScore:          oneModel.BlackScore,
		WhiteScore:          oneModel.WhiteScore,
		NowMoveTime:         config.Conf.Rule.MoveTime,
		NowBlackTime:        oneModel.LeftBlackTime,
		NowWhiteTime:        oneModel.LeftWhiteTime,
		NowBlackByoYomi:     oneModel.LeftBlackByoYomi,
		NowWhiteByoYomi:     oneModel.LeftWhiteByoYomi,
		NowBlackByoYomiTime: oneModel.LeftBlackByoYomiTime,
		NowWhiteByoYomiTime: oneModel.LeftWhiteByoYomiTime,
	}

	cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(request.GameId))
	if !ok {
		e := initService(&InitRequest{GameId: request.GameId})
		if e != nil {
			return e
		}
		cacheData, ok = cache.CachaDataMap.Get(cache.GetKey(request.GameId))
	}
	if ok {
		res.BlackUserOnline = cacheData.BlackUserOnline
		res.WhiteUserOnline = cacheData.WhiteUserOnline
		res.BlackUserEnter = cacheData.BlackUserEnter
		res.WhiteUserEnter = cacheData.WhiteUserEnter
		if cacheData.NowMoveTime != 0 {
			res.NowMoveTime = cacheData.NowMoveTime
		}
		bPlayer := GetPlayer(request.GameId, 1)
		if bPlayer != nil {
			res.NowBlackTime = bPlayer.mainTime
			res.NowBlackByoYomi = bPlayer.byoYomiTime
			res.NowBlackByoYomiTime = bPlayer.byoYomiPeriod
		}
		wPlayer := GetPlayer(request.GameId, 2)
		if wPlayer != nil {
			res.NowWhiteTime = wPlayer.mainTime
			res.NowWhiteByoYomi = wPlayer.byoYomiTime
			res.NowWhiteByoYomiTime = wPlayer.byoYomiPeriod
		}
		if !cacheData.IsEnd {
			if request.UserHash == oneModel.BlackUserHash || request.UserHash == oneModel.WhiteUserHash {
				go SendUserOnline(request.GameId)
			}
		}
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if ok {
		res.BlackCaptured = b.CapturesBy[wq.BLACK]
		res.WhiteCaptured = b.CapturesBy[wq.WHITE]
		res.Turn = b.Player
	}
	return nil
}

func enterService(request *EnterRequest, cacheData *cache.CacheData) error {
	var e error
	ctx, cancel := context.WithCancel(context.Background())
	db := database.GetInstance()
	tx := db.Begin()
	defer func() {
		if e != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			cancel()
		}
	}()
	if cacheData.IsEnd {
		e = errors.New("对弈已结束")
		return e
	}
	if cacheData.CanStartTime > time.Now().Unix() {
		e = errors.New("对弈未到开始时间")
		return e
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok {
		e = errors.New("对弈不存在")
		return e
	}
	switch request.UserHash {
	case cacheData.BlackUserHash:
		if !cacheData.BlackUserEnter {
			cacheData.BlackUserOnline = true
			cacheData.BlackUserEnter = true
			e = model.UpdateGames(tx, &model.Game{Model: gorm.Model{ID: request.GameId}}, &model.Game{BlackUserEnter: true, IsShow: true})
			if e != nil {
				return e
			}
			notifyUserEnter(cacheData.BlackUserId, request.GameId, cacheData.BusinessType, ctx)
			if cacheData.WhiteUserType != model.HUMANUSER {
				e = model.UpdateGames(tx, &model.Game{Model: gorm.Model{ID: request.GameId}}, &model.Game{WhiteUserEnter: true, IsShow: true})
				if e != nil {
					return e
				}
				cacheData.WhiteUserEnter = true
				cacheData.WhiteUserOnline = true
			}
		}
	case cacheData.WhiteUserHash:
		if !cacheData.WhiteUserEnter {
			cacheData.WhiteUserEnter = true
			cacheData.WhiteUserOnline = true
			e = model.UpdateGames(tx, &model.Game{Model: gorm.Model{ID: request.GameId}}, &model.Game{WhiteUserEnter: true, IsShow: true})
			if e != nil {
				return e
			}
			notifyUserEnter(cacheData.WhiteUserId, request.GameId, cacheData.BusinessType, ctx)
			if cacheData.BlackUserType != model.HUMANUSER {
				e = model.UpdateGames(tx, &model.Game{Model: gorm.Model{ID: request.GameId}}, &model.Game{BlackUserEnter: true, IsShow: true})
				if e != nil {
					return e
				}
				cacheData.BlackUserOnline = true
				cacheData.BlackUserEnter = true
				if b.Player == wq.BLACK {
					go aiMove(request.GameId, ctx)
				}
			}
		}
	default:
		e = errors.New("非对弈双方")
		return e
	}
	go startTimeDown(request.GameId, ctx)
	go notifyGameBegin(request.GameId, cacheData.BusinessType, ctx)
	return nil
}

func moveService(request *MoveRequest, cacheData *cache.CacheData, res *MoveResponse) error {
	var e error
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if e == nil {
			cancel()
		}
	}()
	if cacheData.IsEnd {
		e = errors.New("对弈已结束")
		return e
	}
	if !cacheData.TimeDownIsStarted {
		e = errors.New("对弈初始化中")
		return e
	}
	if request.UserHash == cacheData.BlackUserHash && request.C != wq.BLACK {
		e = errors.New("您不能落当前颜色")
		return e
	}
	if request.UserHash == cacheData.WhiteUserHash && request.C != wq.WHITE {
		e = errors.New("您不能落当前颜色")
		return e
	}

	var err error
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok || b == nil {
		e = errors.New("获取对弈信息失败")
		return e
	}
	if b.Paused {
		e = errors.New("对弈已暂停")
		return e
	}
	if b.Player != request.C {
		e = errors.New("落子方不对")
		return e
	}

	cacheData.IsMove = 1
	defer func(c *cache.CacheData) {
		c.IsMove = 0
	}(cacheData)

	err = b.PlayColour(wq.Point(request.X, request.Y), request.C)
	if err != nil {
		e = errors.New("落子失败")
		return e
	}
	setLastMoveCommit(b, request.GameId)
	if b.Player == wq.BLACK {
		b.WContinuePass = 0
	}
	if b.Player == wq.WHITE {
		b.BContinuePass = 0
	}
	b.ContinuePass = 0
	go func(gameId uint, turn wq.Colour) {
		turnCh := GetTurnChan(gameId)
		if turnCh != nil {
			turnCh <- turn
		}
		moveTurnCh := GetMoveTurnChan(gameId)
		if moveTurnCh != nil {
			moveTurnCh <- turn
		}
	}(request.GameId, b.Player)

	go func(gameId uint, x, y int, c wq.Colour, oldTurn wq.Colour, b *wq.Board, ctx context.Context) {
		pointerHash := utils.PointerHash(b.State)
		sendMove(request.GameId, request.X, request.Y, request.C, oldTurn, b.Step, pointerHash, ctx)
	}(request.GameId, request.X, request.Y, request.C, request.C, b, ctx)
	switch b.Player {
	case wq.BLACK:
		if cacheData.BlackUserType != model.HUMANUSER {
			go aiMove(request.GameId, ctx)
		}
	case wq.WHITE:
		if cacheData.WhiteUserType != model.HUMANUSER {
			go aiMove(request.GameId, ctx)
		}
	}

	go func(gameId uint, b *wq.Board) {
		cache.SetGameSgf(gameId, b.GetSGF())
	}(request.GameId, b)
	go func(gameId uint) {
		checkCapture(gameId)
	}(request.GameId)
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		UserId:      request.UserId,
		UserHash:    request.UserHash,
		HandlerTime: time.Now().Unix(),
		Handler:     "move",
		Duration:    time.Now().Unix() - cacheData.LastMoveTime,
	}}, time.Second*5)
	cacheData.LastMoveTime = time.Now().Unix()
	cacheData.LastAreaScoreTime = 0
	cacheData.LastSummationTime = 0
	*res = MoveResponse{
		NextColor: b.Player,
	}
	return nil
}

func passService(request *PassRequest, cacheData *cache.CacheData, res *PassResponse) error {
	var e error
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if e == nil {
			cancel()
		}
	}()
	if !cacheData.TimeDownIsStarted {
		e = errors.New("对弈初始化中")
		return e
	}
	if cacheData.WinCaptured != 0 {
		e = errors.New("吃子模式不允许停一手")
		return e
	}
	if cacheData.IsEnd {
		e = errors.New("对弈已结束")
		return e
	}
	if request.UserHash == cacheData.BlackUserHash && request.C != wq.BLACK {
		e = errors.New("黑方走黑子")
		return e
	}
	if request.UserHash == cacheData.WhiteUserHash && request.C != wq.WHITE {
		e = errors.New("白方走白子")
		return e
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok || b == nil {
		e = errors.New("获取对弈信息失败")
		return e
	}
	if b.Paused {
		e = errors.New("对弈已暂停")
		return e
	}
	if b.Player != request.C {
		e = errors.New("落子方不对")
		return e
	}
	cacheData.IsMove = 1
	defer func(c *cache.CacheData) {
		c.IsMove = 0
	}(cacheData)
	b.Pass()
	setLastMoveCommit(b, request.GameId)
	if b.Step >= cacheData.TerritoryStep {
		if b.Player == wq.BLACK {
			b.WContinuePass++
		}
		if b.Player == wq.WHITE {
			b.BContinuePass++
		}
		b.ContinuePass++
	}
	go func(gameId uint, turn wq.Colour) {
		turnCh := GetTurnChan(gameId)
		if turnCh != nil {
			turnCh <- turn
		}
		moveTurnCh := GetMoveTurnChan(gameId)
		if moveTurnCh != nil {
			moveTurnCh <- turn
		}
	}(request.GameId, b.Player)

	go func(gameId uint, oldTurn wq.Colour, b *wq.Board, ctx context.Context) {
		pointerHash := utils.PointerHash(b.State)
		sendPass(request.GameId, oldTurn, b.Step, pointerHash, ctx)
	}(request.GameId, request.C, b, ctx)
	switch b.Player {
	case wq.BLACK:
		if cacheData.BlackUserType != model.HUMANUSER {
			go aiMove(request.GameId, ctx)
		}
	case wq.WHITE:
		if cacheData.WhiteUserType != model.HUMANUSER {
			go aiMove(request.GameId, ctx)
		}
	}
	go func(gameId uint, b *wq.Board) {
		cache.SetGameSgf(gameId, b.GetSGF())
	}(request.GameId, b)
	go func(gameId uint) {
		checkCapture(gameId)
	}(request.GameId)
	*res = PassResponse{
		NextColor: b.Player,
	}
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		UserId:      request.UserId,
		UserHash:    request.UserHash,
		HandlerTime: time.Now().Unix(),
		Handler:     "pass",
		Duration:    time.Now().Unix() - cacheData.LastMoveTime,
	}}, time.Second*5)
	cacheData.LastMoveTime = time.Now().Unix()
	cacheData.LastAreaScoreTime = 0
	cacheData.LastSummationTime = 0
	go func(b *wq.Board, cacheData *cache.CacheData, r *BaseRequest, ctx context.Context) {
		if b.BContinuePass >= 3 || b.WContinuePass >= 3 || b.ContinuePass >= 2 {
			ticker := time.NewTicker(time.Second * 5)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					var br BaseRequest
					br.GameId = r.GameId
					switch b.Player {
					case wq.BLACK:
						br.UserId = cacheData.BlackUserId
						br.UserHash = cacheData.BlackUserHash
					case wq.WHITE:
						br.UserId = cacheData.WhiteUserId
						br.UserHash = cacheData.WhiteUserHash
					default:
						return
					}
					res := &EndResponse{}
					err := areaScoreService(&InfoRequest{BaseRequest: br}, cacheData, res)
					if err != nil {
						logger.Logger("passService", logger.ERROR, err, fmt.Sprintf("数目失败 err %s", err.Error()))
					}
					return
				case <-ticker.C:
					return
				}
			}
		}
	}(b, cacheData, &request.BaseRequest, ctx)
	return nil
}

func resignService(request *ResignRequest, cacheData *cache.CacheData, res *EndResponse) error {
	var e error
	ctx, cancel := context.WithCancel(context.Background())
	db := database.GetInstance()
	tx := db.Begin()
	defer func() {
		if e != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			cancel()
		}
	}()
	if cacheData.IsEnd {
		e = errors.New("对弈已结束")
		return e
	}
	if request.UserHash != cacheData.BlackUserHash && request.UserHash != cacheData.WhiteUserHash {
		e = errors.New("不是对弈双方")
		return e
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok || b == nil {
		e = errors.New("获取对弈信息失败")
		return e
	}
	b.Paused = true
	var win int
	var winResult string
	switch request.UserHash {
	case cacheData.BlackUserHash:
		win = 2
		winResult = "W+R"
	case cacheData.WhiteUserHash:
		win = 1
		winResult = "B+R"
	}
	shutdownCh := GetShutdownChan(request.GameId)
	if shutdownCh != nil {
		select {
		case shutdownCh <- struct{}{}:
		case <-time.After(2 * time.Second):
			logger.Logger("resignService", logger.ERROR, nil, fmt.Sprintf("gameId:%d, shutdownCh is closed", request.GameId))
		}
	}
	moveTickerEndCh := GetMoveTickEndChan(request.GameId)
	if moveTickerEndCh != nil {
		select {
		case moveTickerEndCh <- struct{}{}:
		case <-time.After(2 * time.Second):
			logger.Logger("resignService", logger.ERROR, nil, fmt.Sprintf("gameId:%d, moveTickerEndCh is closed", request.GameId))
		}
	}
	e = model.UpdateGames(tx,
		&model.Game{Model: gorm.Model{ID: request.GameId}},
		&model.Game{IsEnd: true, EndTime: time.Now().Unix(), WinResult: winResult, Win: win, SGF: b.GetSGF(), Step: b.Step, BlackCaptured: b.CapturesBy[wq.BLACK], WhiteCaptured: b.CapturesBy[wq.WHITE], IsShow: true})
	if e != nil {
		return e
	}
	sendEnd(request.GameId, win, winResult, ctx, "resignService")
	notifyGameEnd(request.GameId, win, winResult, cacheData, b, ctx)
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		HandlerTime: time.Now().Unix(),
		Handler:     "end",
		Content:     b.GetCommitSGF(),
	}}, time.Second*5)
	cacheData.IsEnd = true
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		UserId:      request.UserId,
		UserHash:    request.UserHash,
		HandlerTime: time.Now().Unix(),
		Handler:     "resign",
	}}, time.Second*5)
	*res = EndResponse{
		Win:       win,
		WinResult: winResult,
	}
	return nil
}

func areaScoreService(request *InfoRequest, cacheData *cache.CacheData, reply *EndResponse) error {
	var e error
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if e == nil {
			cancel()
		}
	}()
	if cacheData.IsEnd {
		e = errors.New("对弈已结束")
		return e
	}
	if request.UserHash != cacheData.BlackUserHash && request.UserHash != cacheData.WhiteUserHash {
		e = errors.New("不是对弈双方")
		return e
	}

	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok || b == nil {
		e = errors.New("获取对弈信息失败")
		return e
	}
	if b.Paused {
		e = errors.New("对弈已暂停")
		return e
	}
	switch request.UserHash {
	case cacheData.BlackUserHash:
		if b.Player != wq.BLACK {
			e = errors.New("当前轮到白方")
			return e
		}
	case cacheData.WhiteUserHash:
		if b.Player != wq.WHITE {
			e = errors.New("当前轮到黑方")
			return e
		}
	}
	b.Paused = true
	bScore, wScore, endScore, controversyCount, ownership, err := areaScore(request.GameId, b)
	if err != nil {
		b.Paused = false
		e = errors.New("数目失败")
		return e
	}
	db := database.GetInstance()
	tx := db.Begin()
	b.BScore = bScore
	b.WScore = wScore
	b.ControversyCount = controversyCount
	b.OwnerShip = *ownership
	win, winResult, err := scoreEnd(ctx, tx, endScore, request.GameId, b, cacheData)
	if err != nil {
		e = err
		tx.Rollback()
		return e
	}
	cacheData.IsEnd = true
	tx.Commit()
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		UserId:      request.UserId,
		UserHash:    request.UserHash,
		HandlerTime: time.Now().Unix(),
		Handler:     "areaScore",
	}}, time.Second*5)
	reply.Win = win
	reply.WinResult = winResult
	return nil
}

func applyForAreaScoreService(request *InfoRequest, cacheData *cache.CacheData) error {
	var e error
	if cacheData.IsEnd {
		e = errors.New("对弈已结束")
		return e
	}
	if request.UserHash != cacheData.BlackUserHash && request.UserHash != cacheData.WhiteUserHash {
		e = errors.New("不是对弈双方")
		return e
	}
	if time.Now().Unix()-cacheData.LastAreaScoreTime <= config.Conf.Rule.AreaScoreTimeInterval {
		e = errors.New("申请太频繁")
		return e
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok || b == nil {
		e = errors.New("获取对弈信息失败")
		return e
	}
	if b.Step < cacheData.TerritoryStep {
		e = errors.New("步数不够")
		return e
	}
	if b.Paused {
		e = errors.New("对弈已暂停")
		return e
	}
	switch request.UserHash {
	case cacheData.BlackUserHash:
		if b.Player != wq.BLACK {
			e = errors.New("当前轮到白方")
			return e
		}
	case cacheData.WhiteUserHash:
		if b.Player != wq.WHITE {
			e = errors.New("当前轮到黑方")
			return e
		}
	}
	SetAreaScoreTimeDownApplyUser(request.GameId, request.UserHash)
	b.Paused = true
	var opponentUserId string
	var opponentUserHash string
	switch request.UserHash {
	case cacheData.BlackUserHash:
		opponentUserId = cacheData.WhiteUserId
		opponentUserHash = cacheData.WhiteUserHash
		if cacheData.WhiteUserType != model.HUMANUSER {
			go FakeUserAgreeAreaScore(&InfoRequest{
				BaseRequest: BaseRequest{
					GameId:   request.GameId,
					UserId:   cacheData.WhiteUserId,
					UserHash: cacheData.WhiteUserHash,
				},
			})
		}
	case cacheData.WhiteUserHash:
		opponentUserId = cacheData.BlackUserId
		opponentUserHash = cacheData.BlackUserHash
		if cacheData.BlackUserType != model.HUMANUSER {
			go FakeUserAgreeAreaScore(&InfoRequest{
				BaseRequest: BaseRequest{
					GameId:   request.GameId,
					UserId:   cacheData.BlackUserId,
					UserHash: cacheData.BlackUserHash,
				},
			})
		}
	}

	go func(gameId uint, userId string, userHash string, opponentUserId string, opponentUserHash string) {
		var agreeCh chan struct{}
		agreeCh = make(chan struct{})
		var rejectCh chan struct{}
		rejectCh = make(chan struct{})
		SetAreaScoreTimeDownAgreeChan(request.GameId, agreeCh)
		SetAreaScoreTimeDownRejectChan(request.GameId, rejectCh)
		AreaScoreTimeDown(gameId, userId, userHash, opponentUserId, opponentUserHash, agreeCh, rejectCh)
	}(request.GameId, request.UserId, request.UserHash, opponentUserId, opponentUserHash)
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		UserId:      request.UserId,
		UserHash:    request.UserHash,
		Handler:     "applyForAreaScoreService",
		HandlerTime: time.Now().Unix(),
	}}, time.Second*5)
	return nil
}

func agreeAreaScoreService(request *InfoRequest, cacheData *cache.CacheData, reply *EndResponse) error {
	var e error
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if e == nil {
			cancel()
		}
	}()
	if cacheData.IsEnd {
		e = errors.New("对弈已结束")
		return e
	}
	if request.UserHash != cacheData.BlackUserHash && request.UserHash != cacheData.WhiteUserHash {
		e = errors.New("不是对弈双方")
		return e
	}
	applyUserHash := GetAreaScoreTimeDownApplyUser(request.GameId)
	if applyUserHash == "" {
		e = errors.New("未申请数目")
		return e
	}
	var opponentUserId string
	switch applyUserHash {
	case cacheData.BlackUserHash:
		opponentUserId = cacheData.WhiteUserId
	case cacheData.WhiteUserHash:
		opponentUserId = cacheData.BlackUserId
	}
	if opponentUserId != request.UserId {
		e = errors.New("该用户不能同意数目")
		return e
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok || b == nil {
		e = errors.New("获取对弈信息失败")
		return e
	}
	if !b.Paused {
		e = errors.New("对弈未暂停")
		return e
	}

	agreeCh := GetAreaScoreTimeDownAgreeChan(request.GameId)
	if agreeCh != nil {
		select {
		case agreeCh <- struct{}{}:
		case <-time.After(2 * time.Second):
			logger.Logger("agreeAreaScoreService", logger.ERROR, nil, fmt.Sprintf("gameId:%d, agreeCh is closed", request.GameId))
		}
	}
	bScore, wScore, endScore, controversyCount, ownership, err := areaScore(request.GameId, b)
	if err != nil {
		b.Paused = false
		e = errors.New("数目失败")
		return e
	}
	db := database.GetInstance()
	tx := db.Begin()
	b.BScore = bScore
	b.WScore = wScore
	b.ControversyCount = controversyCount
	b.OwnerShip = *ownership
	win, winResult, err := scoreEnd(ctx, tx, endScore, request.GameId, b, cacheData)
	if err != nil {
		e = err
		tx.Rollback()
		return e
	}
	cacheData.IsEnd = true
	tx.Commit()
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		UserId:      request.UserId,
		UserHash:    request.UserHash,
		Handler:     "agreeAreaScoreService",
		HandlerTime: time.Now().Unix(),
	}}, time.Second*5)
	reply.Win = win
	reply.WinResult = winResult
	return nil
}

func rejectAreaScoreService(request *InfoRequest, cacheData *cache.CacheData) error {
	var e error
	if cacheData.IsEnd {
		e = errors.New("对弈已结束")
		return e
	}
	if request.UserHash != cacheData.BlackUserHash && request.UserHash != cacheData.WhiteUserHash {
		e = errors.New("不是对弈双方")
		return e
	}
	applyUserHash := GetAreaScoreTimeDownApplyUser(request.GameId)
	if applyUserHash == "" {
		e = errors.New("未申请数目")
		return e
	}
	var opponentUserId string
	switch applyUserHash {
	case cacheData.BlackUserHash:
		opponentUserId = cacheData.WhiteUserId
	case cacheData.WhiteUserHash:
		opponentUserId = cacheData.BlackUserId
	}
	if opponentUserId != request.UserId {
		e = errors.New("该用户不能拒绝数目")
		return e
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok || b == nil {
		e = errors.New("获取对弈信息失败")
		return e
	}
	if !b.Paused {
		e = errors.New("对弈未暂停")
		return e
	}
	rejectCh := GetAreaScoreTimeDownRejectChan(request.GameId)
	if rejectCh != nil {
		select {
		case rejectCh <- struct{}{}:
		case <-time.After(2 * time.Second):
			logger.Logger("rejectAreaScoreService", logger.ERROR, nil, fmt.Sprintf("gameId:%d, rejectCh is closed", request.GameId))
		}
	}
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		UserId:      request.UserId,
		UserHash:    request.UserHash,
		Handler:     "rejectAreaScoreService",
		HandlerTime: time.Now().Unix(),
	}}, time.Second*5)
	return nil
}

func drawService(request *InfoRequest, cacheData *cache.CacheData) error {
	var e error
	db := database.GetInstance()
	tx := db.Begin()
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if e != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			cancel()
		}
	}()
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok || b == nil {
		e = errors.New("获取对弈信息失败")
		return e
	}
	b.Paused = true
	shutdownCh := GetShutdownChan(request.GameId)
	if shutdownCh != nil {
		select {
		case shutdownCh <- struct{}{}:
		case <-time.After(2 * time.Second):
			logger.Logger("drawService", logger.ERROR, nil, fmt.Sprintf("gameId:%d, shutdownCh is closed", request.GameId))
		}
	}
	var win int
	var winResult string
	// 和棋
	win = 3
	winResult = "Draw"
	e = model.UpdateGames(tx,
		&model.Game{Model: gorm.Model{ID: request.GameId}},
		&model.Game{IsEnd: true, EndTime: time.Now().Unix(), WinResult: winResult, Win: win, SGF: b.GetSGF(), Step: b.Step, BlackCaptured: b.CapturesBy[wq.BLACK], WhiteCaptured: b.CapturesBy[wq.WHITE], IsShow: true})
	if e != nil {
		return e
	}
	sendEnd(request.GameId, win, winResult, ctx, "drawService")
	notifyGameEnd(request.GameId, win, winResult, cacheData, b, ctx)
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		HandlerTime: time.Now().Unix(),
		Handler:     "end",
		Content:     b.GetCommitSGF(),
	}}, time.Second*5)
	return nil
}

func callEndService(request *CallEndRequest) error {
	var e error
	db := database.GetInstance()
	tx := db.Begin()
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if e != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			cancel()
		}
	}()
	g := &model.Game{
		SGF:       request.SGF,
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
		Win:       request.Win,
		WinResult: request.WinResult,
		IsShow:    true,
	}
	if request.SGF != "" {
		n, _, err := wq.Load(request.SGF)
		if err != nil {
			e = err
			return e
		}
		b := n.Board()
		g.SGF = request.SGF
		g.Step = b.Step
		g.BlackCaptured = b.CapturesBy[wq.BLACK]
		g.WhiteCaptured = b.CapturesBy[wq.WHITE]
	}
	if g.StartTime != 0 {
		g.IsStart = true
	}
	if g.EndTime != 0 {
		g.IsEnd = true
	}

	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if ok && b != nil {
		b.Paused = true
		shutdownCh := GetShutdownChan(request.GameId)
		if shutdownCh != nil {
			select {
			case shutdownCh <- struct{}{}:
			case <-time.After(time.Second):
				logger.Logger("callEndService", logger.ERROR, nil, fmt.Sprintf("gameId:%d, shutdownCh is closed", request.GameId))
			}
		}
		moveTickerEndCh := GetMoveTickEndChan(request.GameId)
		if moveTickerEndCh != nil {
			select {
			case moveTickerEndCh <- struct{}{}:
			case <-time.After(time.Second):
				logger.Logger("callEndService", logger.ERROR, nil, fmt.Sprintf("gameId:%d, moveTickerEndCh is closed", request.GameId))
			}
		}
		if g.SGF != "" {
			cache.SetGameSgf(request.GameId, g.SGF)
		}
	}
	cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(request.GameId))
	if ok && cacheData != nil {
		cacheData.IsEnd = true
	}
	sendEnd(request.GameId, request.Win, request.WinResult, ctx, "callEndService")
	notifyGameEnd(request.GameId, request.Win, request.WinResult, cacheData, b, ctx)
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		HandlerTime: time.Now().Unix(),
		Handler:     "end",
		Content:     b.GetCommitSGF(),
	}}, time.Second*5)
	e = model.UpdateGames(tx, &model.Game{Model: gorm.Model{ID: request.GameId}}, g)
	if e != nil {
		return e
	}
	return nil
}

func applyForSummationService(request *InfoRequest, cacheData *cache.CacheData) error {
	var e error
	if cacheData.IsEnd {
		e = errors.New("对弈已结束")
		return e
	}
	if request.UserHash != cacheData.BlackUserHash && request.UserHash != cacheData.WhiteUserHash {
		e = errors.New("不是对弈双方")
		return e
	}
	if cacheData.SummationCount >= config.Conf.Rule.SummationCount {
		e = errors.New("申请超过限制")
		return e
	}
	if time.Now().Unix()-cacheData.LastSummationTime <= config.Conf.Rule.SummationTimeInterval {
		e = errors.New("申请太频繁")
		return e
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok || b == nil {
		e = errors.New("获取对弈信息失败")
		return e
	}
	if b.Paused {
		e = errors.New("对弈已暂停")
		return e
	}
	switch request.UserHash {
	case cacheData.BlackUserHash:
		if b.Player != wq.BLACK {
			e = errors.New("当前轮到白方")
			return e
		}
	case cacheData.WhiteUserHash:
		if b.Player != wq.WHITE {
			e = errors.New("当前轮到黑方")
			return e
		}
	}
	cacheData.SummationCount++
	SetSummationTimeDownApplyUser(request.GameId, request.UserHash)
	b.Paused = true
	var opponentUserId string
	var opponentUserHash string
	switch request.UserHash {
	case cacheData.BlackUserHash:
		opponentUserId = cacheData.WhiteUserId
		opponentUserHash = cacheData.WhiteUserHash
		if cacheData.WhiteUserType != model.HUMANUSER {
			go FakeUserAgreeSummation(&InfoRequest{
				BaseRequest: BaseRequest{
					GameId:   request.GameId,
					UserId:   cacheData.WhiteUserId,
					UserHash: cacheData.WhiteUserHash,
				},
			})
		}
	case cacheData.WhiteUserHash:
		opponentUserId = cacheData.BlackUserId
		opponentUserHash = cacheData.BlackUserHash
		if cacheData.BlackUserType != model.HUMANUSER {
			go FakeUserAgreeSummation(&InfoRequest{
				BaseRequest: BaseRequest{
					GameId:   request.GameId,
					UserId:   cacheData.BlackUserId,
					UserHash: cacheData.BlackUserHash,
				},
			})
		}
	}

	go func(gameId uint, userId string, userHash string, opponentUserId string, opponentUserHash string) {
		var agreeCh chan struct{}
		agreeCh = make(chan struct{})
		var rejectCh chan struct{}
		rejectCh = make(chan struct{})
		SetSummationTimeDownAgreeChan(request.GameId, agreeCh)
		SetSummationTimeDownRejectChan(request.GameId, rejectCh)
		SummationTimeDown(gameId, userId, userHash, opponentUserId, opponentUserHash, agreeCh, rejectCh)
	}(request.GameId, request.UserId, request.UserHash, opponentUserId, opponentUserHash)
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		UserId:      request.UserId,
		UserHash:    request.UserHash,
		Handler:     "applyForSummationService",
		HandlerTime: time.Now().Unix(),
	}}, time.Second*5)
	return nil
}

func agreeSummationService(request *InfoRequest, cacheData *cache.CacheData, res *EndResponse) error {
	var e error
	db := database.GetInstance()
	tx := db.Begin()
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if e != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			cancel()
		}
	}()
	if cacheData.IsEnd {
		e = errors.New("对弈已结束")
		return e
	}
	if request.UserHash != cacheData.BlackUserHash && request.UserHash != cacheData.WhiteUserHash {
		e = errors.New("不是对弈双方")
		return e
	}
	applyUserHash := GetSummationTimeDownApplyUser(request.GameId)
	if applyUserHash == "" {
		e = errors.New("未申请和棋")
		return e
	}
	var opponentUserId string
	switch applyUserHash {
	case cacheData.BlackUserHash:
		opponentUserId = cacheData.WhiteUserId
	case cacheData.WhiteUserHash:
		opponentUserId = cacheData.BlackUserId
	}
	if opponentUserId != request.UserId {
		e = errors.New("该用户不能同意数目")
		return e
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok || b == nil {
		e = errors.New("获取对弈信息失败")
		return e
	}
	if !b.Paused {
		e = errors.New("对弈未暂停")
		return e
	}
	agreeCh := GetSummationTimeDownAgreeChan(request.GameId)
	if agreeCh != nil {
		select {
		case agreeCh <- struct{}{}:
		case <-time.After(time.Second):
			logger.Logger("agreeSummationService", logger.ERROR, nil, fmt.Sprintf("gameId:%d, agreeCh is closed", request.GameId))
		}
	}
	win := 3
	winResult := "Draw"
	shutdownCh := GetShutdownChan(request.GameId)
	if shutdownCh != nil {
		select {
		case shutdownCh <- struct{}{}:
		case <-time.After(time.Second):
			logger.Logger("agreeSummationService", logger.ERROR, nil, fmt.Sprintf("gameId:%d, shutdownCh is closed", request.GameId))
		}
	}
	moveTickerEndCh := GetMoveTickEndChan(request.GameId)
	if moveTickerEndCh != nil {
		select {
		case moveTickerEndCh <- struct{}{}:
		case <-time.After(time.Second):
			logger.Logger("agreeSummationService", logger.ERROR, nil, fmt.Sprintf("gameId:%d, moveTickerEndCh is closed", request.GameId))
		}
	}
	e = model.UpdateGames(tx,
		&model.Game{Model: gorm.Model{ID: request.GameId}},
		&model.Game{IsEnd: true, EndTime: time.Now().Unix(), WinResult: winResult, Win: win, SGF: b.GetSGF(), Step: b.Step, BlackCaptured: b.CapturesBy[wq.BLACK], WhiteCaptured: b.CapturesBy[wq.WHITE], IsShow: true})
	if e != nil {
		return e
	}
	sendEnd(request.GameId, win, winResult, ctx, "agreeSummationService")
	notifyGameEnd(request.GameId, win, winResult, cacheData, b, ctx)

	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		HandlerTime: time.Now().Unix(),
		Handler:     "end",
		Content:     b.GetCommitSGF(),
	}}, time.Second*5)

	cacheData.IsEnd = true
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		UserId:      request.UserId,
		UserHash:    request.UserHash,
		Handler:     "agreeSummationService",
		HandlerTime: time.Now().Unix(),
	}}, time.Second*5)
	*res = EndResponse{
		Win:       win,
		WinResult: winResult,
	}
	return nil
}

func rejectSummationService(request *InfoRequest, cacheData *cache.CacheData) error {
	var e error
	if cacheData.IsEnd {
		e = errors.New("对弈已结束")
		return e
	}
	if request.UserHash != cacheData.BlackUserHash && request.UserHash != cacheData.WhiteUserHash {
		e = errors.New("不是对弈双方")
		return e
	}
	applyUserHash := GetSummationTimeDownApplyUser(request.GameId)
	if applyUserHash == "" {
		e = errors.New("未申请和棋")
		return e
	}
	var opponentUserId string
	switch applyUserHash {
	case cacheData.BlackUserHash:
		opponentUserId = cacheData.WhiteUserId
	case cacheData.WhiteUserHash:
		opponentUserId = cacheData.BlackUserId
	}
	if opponentUserId != request.UserId {
		e = errors.New("该用户不能拒绝数目")
		return e
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok || b == nil {
		e = errors.New("获取对弈信息失败")
		return e
	}
	if !b.Paused {
		e = errors.New("对弈未暂停")
		return e
	}
	rejectCh := GetSummationTimeDownRejectChan(request.GameId)
	if rejectCh != nil {
		select {
		case rejectCh <- struct{}{}:
		case <-time.After(time.Second):
			logger.Logger("rejectSummationService", logger.ERROR, nil, fmt.Sprintf("gameId:%d, rejectCh is closed", request.GameId))
		}
	}
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      request.GameId,
		UserId:      request.UserId,
		UserHash:    request.UserHash,
		Handler:     "rejectSummationService",
		HandlerTime: time.Now().Unix(),
	}}, time.Second*5)
	return nil
}

func canPlayService(request *InfoRequest, cacheData *cache.CacheData, res *CanPlayResponse) error {
	var e error
	if cacheData.IsEnd {
		e = errors.New("对弈已结束")
		return e
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(request.GameId))
	if !ok || b == nil {
		e = errors.New("获取对弈信息失败")
		return e
	}
	if b.Paused {
		e = errors.New("对弈已暂停")
		return e
	}
	if request.UserHash == cacheData.BlackUserHash && b.Player != wq.BLACK {
		*res = CanPlayResponse{NextColor: 0}
		return nil
	}
	if request.UserHash == cacheData.WhiteUserHash && b.Player != wq.WHITE {
		*res = CanPlayResponse{NextColor: 0}
		return nil
	}
	*res = CanPlayResponse{NextColor: b.Player}
	return nil
}

func forceReloadService(request *ForceReloadSgfRequest) error {
	var e error
	newGame, _, err := wq.Load(request.SGF)
	if err != nil {
		e = err
		return e
	}
	b := newGame.Board()
	wq.WQDB.Set(wq.GameIdKey(request.GameId), b)
	d := notifyStruct.ForceReloadSGFMsg{
		GameId: request.GameId,
		SGF:    request.SGF,
	}
	ws.GroupPublishApiChan <- &ws.GroupPublishData{
		GroupId: fmt.Sprintf("game:%d", request.GameId),
		Message: notifyStruct.MakeWsForceReloadSGFMsg(&d),
	}
	cache.SetGameSgf(request.GameId, request.SGF)
	return nil
}

func setLastMoveCommit(b *wq.Board, gameId uint) {
	b.Move[len(b.Move)-1].Commit.TimeStamp = time.Now().Unix()
	b.Move[len(b.Move)-1].Commit.HasCommit = true
	bPlayer := GetPlayer(gameId, 1)
	if bPlayer != nil {
		b.Move[len(b.Move)-1].Commit.BTime = bPlayer.mainTime
		b.Move[len(b.Move)-1].Commit.BByo = bPlayer.byoYomiTime
		b.Move[len(b.Move)-1].Commit.BByoT = bPlayer.byoYomiPeriod
	}
	wPlayer := GetPlayer(gameId, 2)
	if wPlayer != nil {
		b.Move[len(b.Move)-1].Commit.WTime = wPlayer.mainTime
		b.Move[len(b.Move)-1].Commit.WByo = wPlayer.byoYomiTime
		b.Move[len(b.Move)-1].Commit.WByoT = wPlayer.byoYomiPeriod
	}
}
