package play

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"higo-game-node/api/ws"
	"higo-game-node/cache"
	"higo-game-node/database"
	"higo-game-node/handlerHistory"
	"higo-game-node/logger"
	"higo-game-node/model"
	"higo-game-node/notifyStruct"
	"higo-game-node/redisUtils"
	"higo-game-node/sharedMap"
	"higo-game-node/wq"
	"math/rand"
	"sync"
	"time"
)

var turnChanMap = sharedMap.New[chan wq.Colour]()

func GetTurnChan(gameId uint) chan wq.Colour {
	value, ok := turnChanMap.Get(fmt.Sprintf("%d", gameId))
	if !ok {
		return nil
	}
	return value
}

func SetTurnChan(gameId uint, turnCh chan wq.Colour) {
	turnChanMap.Set(fmt.Sprintf("%d", gameId), turnCh)
}

func DeleteTurnChan(gameId uint) {
	turnChanMap.Remove(fmt.Sprintf("%d", gameId))
}

var shutdownChanMap = sharedMap.New[chan struct{}]()

func GetShutdownChan(gameId uint) chan struct{} {
	value, ok := shutdownChanMap.Get(fmt.Sprintf("%d", gameId))
	if !ok {
		return nil
	}
	return value
}
func SetShutdownChan(gameId uint, shutdownCh chan struct{}) {
	shutdownChanMap.Set(fmt.Sprintf("%d", gameId), shutdownCh)
}

func DeleteShutdownChan(gameId uint) {
	shutdownChanMap.Remove(fmt.Sprintf("%d", gameId))
}

type Player struct {
	gameId             uint
	color              wq.Colour
	currentColor       wq.Colour
	mainTime           int // 秒
	byoYomiTime        int // 秒
	defaultByoYomiTime int // 秒
	byoYomiPeriod      int // 秒
	sendStartByomi     bool
}

var playerMap = sharedMap.New[*Player]()

func GetPlayer(gameId uint, color wq.Colour) *Player {
	value, ok := playerMap.Get(fmt.Sprintf("%d:%d", gameId, color))
	if !ok {
		return nil
	}
	return value
}

func SetPlayer(gameId uint, color wq.Colour, player *Player) {
	playerMap.Set(fmt.Sprintf("%d:%d", gameId, color), player)
}

func DeletePlayer(gameId uint, color wq.Colour) {
	playerMap.Remove(fmt.Sprintf("%d:%d", gameId, color))
}

func (p *Player) StartTimer(callEnd chan struct{}, notifyCh chan struct{}, otherPlayer *Player) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("Player StartTimer", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", p.gameId, err))
		}
	}()
	count := 0
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	b, ok := wq.WQDB.Get(wq.GameIdKey(p.gameId))
	if !ok {
		return
	}
	for {
		select {
		case <-callEnd:
			return
		case <-ticker.C:
			if p == nil || otherPlayer == nil {
				return
			}
			if p.currentColor == p.color && !b.Paused && b.Player == p.color {
				d := notifyStruct.GameTimeMsg{
					GameId: p.gameId,
					C:      p.currentColor,
				}
				if count%10 == 0 {
					switch p.color {
					case wq.BLACK:
						d.BlackMainTime = p.mainTime
						d.BlackByoYomiTime = p.byoYomiTime
						d.BlackByoYomi = p.byoYomiPeriod
						d.WhiteMainTime = otherPlayer.mainTime
						d.WhiteByoYomiTime = otherPlayer.byoYomiTime
						d.WhiteByoYomi = otherPlayer.byoYomiPeriod
					case wq.WHITE:
						d.WhiteMainTime = p.mainTime
						d.WhiteByoYomiTime = p.byoYomiTime
						d.WhiteByoYomi = p.byoYomiPeriod
						d.BlackMainTime = otherPlayer.mainTime
						d.BlackByoYomiTime = otherPlayer.byoYomiTime
						d.BlackByoYomi = otherPlayer.byoYomiPeriod
					default:
						continue
					}
					ws.GroupPublishApiChan <- &ws.GroupPublishData{
						GroupId: fmt.Sprintf("game:%d", p.gameId),
						Message: notifyStruct.MakeWsGameTimeMsg(&d),
						Ctx:     nil,
					}
				}
				if p.mainTime > 0 {
					p.mainTime--
					if p.mainTime <= 30 {
						count = 0
					} else {
						count++
					}
				} else if p.mainTime == 0 {
					if !p.sendStartByomi {
						p.sendStartByomi = true
						x := notifyStruct.StartByomiMsg{
							GameId: p.gameId,
							C:      p.color,
						}
						ws.GroupPublishApiChan <- &ws.GroupPublishData{
							GroupId: fmt.Sprintf("game:%d", p.gameId),
							Message: notifyStruct.MakeWsStartByomiMsg(&x),
						}
					}
					count = 0
					if p.byoYomiTime > 0 {
						p.byoYomiTime--
					} else if p.byoYomiTime == 0 {
						if p.byoYomiPeriod > 1 {
							p.byoYomiPeriod--
							p.byoYomiTime = p.defaultByoYomiTime
						} else {
							notifyCh <- struct{}{}
							return
						}
					} else {
						notifyCh <- struct{}{}
						return
					}
				}
			} else {
				count = 0
			}
		}
	}
}

type gameTime struct {
	BlackTime               int
	WhiteTime               int
	BlackByoYomi            int
	WhiteByoYomi            int
	BlackByoYomiTime        int
	BlackDefaultByoYomiTime int
	WhiteByoYomiTime        int
	WhiteDefaultByoYomiTime int
}

func StartTimeDown(gameId uint, timeData *gameTime, doneFunc func()) {
	once := sync.Once{}

	defer func() {
		if err := recover(); err != nil {
			once.Do(doneFunc)
			logger.Logger("StartTimeDown", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	var turnCh chan wq.Colour
	turnCh = make(chan wq.Colour)
	var shutdownCh chan struct{}
	shutdownCh = make(chan struct{})
	SetTurnChan(gameId, turnCh)
	SetShutdownChan(gameId, shutdownCh)
	defer DeleteTurnChan(gameId)
	defer DeleteShutdownChan(gameId)

	if (timeData.BlackTime == 0 && timeData.BlackByoYomi == 0 && timeData.BlackByoYomiTime == 0) ||
		(timeData.WhiteTime == 0 && timeData.WhiteByoYomi == 0 && timeData.WhiteByoYomiTime == 0) {
		once.Do(doneFunc)
		return
	}

	blackPlayer := &Player{
		gameId:             gameId,
		color:              wq.BLACK,
		currentColor:       wq.BLACK,
		mainTime:           timeData.BlackTime,
		byoYomiTime:        timeData.BlackByoYomiTime,
		defaultByoYomiTime: timeData.BlackDefaultByoYomiTime,
		byoYomiPeriod:      timeData.BlackByoYomi,
	}
	SetPlayer(gameId, wq.BLACK, blackPlayer)
	defer DeletePlayer(gameId, wq.BLACK)
	whitePlayer := &Player{
		gameId:             gameId,
		color:              wq.WHITE,
		currentColor:       wq.BLACK,
		mainTime:           timeData.WhiteTime,
		byoYomiTime:        timeData.WhiteByoYomiTime,
		defaultByoYomiTime: timeData.WhiteDefaultByoYomiTime,
		byoYomiPeriod:      timeData.WhiteByoYomi,
	}
	SetPlayer(gameId, wq.WHITE, whitePlayer)
	defer DeletePlayer(gameId, wq.WHITE)
	var bTimeEndCh chan struct{}
	bTimeEndCh = make(chan struct{})
	var bNotifyCh chan struct{}
	bNotifyCh = make(chan struct{})
	var wTimeEndCh chan struct{}
	wTimeEndCh = make(chan struct{})
	var wNotifyCh chan struct{}
	wNotifyCh = make(chan struct{})
	go blackPlayer.StartTimer(bTimeEndCh, bNotifyCh, whitePlayer)
	go whitePlayer.StartTimer(wTimeEndCh, wNotifyCh, blackPlayer)

	for {
		once.Do(doneFunc)
		select {
		case color := <-turnCh:
			blackPlayer.currentColor = color
			whitePlayer.currentColor = color
			blackPlayer.byoYomiTime = timeData.BlackDefaultByoYomiTime
			whitePlayer.byoYomiTime = timeData.WhiteDefaultByoYomiTime
			SetGamePlayerTimeToRedis(gameId, wq.BLACK, blackPlayer.mainTime, blackPlayer.byoYomiPeriod, blackPlayer.byoYomiTime)
			SetGamePlayerTimeToRedis(gameId, wq.WHITE, whitePlayer.mainTime, whitePlayer.byoYomiPeriod, whitePlayer.byoYomiTime)
		case <-shutdownCh:
			select {
			case bTimeEndCh <- struct{}{}:
			case <-time.After(time.Second):
				logger.Logger("StartTimeDown", logger.ERROR, nil, "shutdownCh bTimeEndCh error")
			}
			select {
			case wTimeEndCh <- struct{}{}:
			case <-time.After(time.Second):
				logger.Logger("StartTimeDown", logger.ERROR, nil, "shutdownCh wTimeEndCh error")
			}
			db := database.GetInstance()
			tx := db.Begin()
			e := model.UpdateGames(tx, &model.Game{Model: gorm.Model{ID: gameId}},
				&model.Game{
					LeftBlackTime:        blackPlayer.mainTime,
					LeftBlackByoYomi:     blackPlayer.byoYomiPeriod,
					LeftBlackByoYomiTime: blackPlayer.byoYomiTime,
					LeftWhiteTime:        whitePlayer.mainTime,
					LeftWhiteByoYomi:     whitePlayer.byoYomiPeriod,
					LeftWhiteByoYomiTime: whitePlayer.byoYomiTime,
					IsShow:               true,
				}, func(db *gorm.DB) *gorm.DB {
					return db.Select("left_black_time", "left_black_byo_yomi", "left_black_byo_yomi_time", "left_white_time", "left_white_byo_yomi", "left_white_byo_yomi_time", "is_show")
				})
			if e != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
			return
		case <-bNotifyCh:
			select {
			case wTimeEndCh <- struct{}{}:
			case <-time.After(time.Second):
				logger.Logger("StartTimeDown", logger.ERROR, nil, "bNotifyCh wTimeEndCh error")
			}
			g := model.Game{
				IsEnd:                true,
				EndTime:              time.Now().Unix(),
				LeftBlackTime:        blackPlayer.mainTime,
				LeftBlackByoYomi:     blackPlayer.byoYomiPeriod,
				LeftBlackByoYomiTime: blackPlayer.byoYomiTime,
				LeftWhiteTime:        whitePlayer.mainTime,
				LeftWhiteByoYomi:     whitePlayer.byoYomiPeriod,
				LeftWhiteByoYomiTime: whitePlayer.byoYomiTime,
				Win:                  2,
				WinResult:            "W+T",
			}
			b, ok := wq.WQDB.Get(wq.GameIdKey(gameId))
			if ok {
				g.SGF = b.GetSGF()
				g.Step = b.Step
				g.BlackCaptured = b.CapturesBy[wq.BLACK]
				g.WhiteCaptured = b.CapturesBy[wq.WHITE]
			}
			saveTimeEnd(&g, gameId, b, "bTimeEnd")
		case <-wNotifyCh:
			select {
			case bTimeEndCh <- struct{}{}:
			case <-time.After(time.Second):
				logger.Logger("StartTimeDown", logger.ERROR, nil, "bNotifyCh bTimeEndCh error")
			}
			g := model.Game{
				IsEnd:                true,
				EndTime:              time.Now().Unix(),
				LeftBlackTime:        blackPlayer.mainTime,
				LeftBlackByoYomi:     blackPlayer.byoYomiPeriod,
				LeftBlackByoYomiTime: blackPlayer.byoYomiTime,
				LeftWhiteTime:        whitePlayer.mainTime,
				LeftWhiteByoYomi:     whitePlayer.byoYomiPeriod,
				LeftWhiteByoYomiTime: whitePlayer.byoYomiTime,
				Win:                  1,
				WinResult:            "B+T",
				IsShow:               true,
			}
			b, ok := wq.WQDB.Get(wq.GameIdKey(gameId))
			if ok {
				g.SGF = b.GetSGF()
				g.Step = b.Step
				g.BlackCaptured = b.CapturesBy[wq.BLACK]
				g.WhiteCaptured = b.CapturesBy[wq.WHITE]
			}
			saveTimeEnd(&g, gameId, b, "wTimeEnd")
			return
		}
	}
}

func saveTimeEnd(g *model.Game, gameId uint, b *wq.Board, from string) {
	db := database.GetInstance()
	tx := db.Begin()
	e := model.UpdateGames(tx, &model.Game{Model: gorm.Model{ID: gameId}}, g, func(db *gorm.DB) *gorm.DB {
		return db.Select("is_end",
			"end_time",
			"left_black_time",
			"left_black_byo_yomi",
			"left_black_byo_yomi_time",
			"left_white_time",
			"left_white_byo_yomi",
			"left_white_byo_yomi_time",
			"sgf",
			"step",
			"black_captured",
			"white_captured",
			"win",
			"win_result", "is_show")
	})
	if e != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	sendEnd(gameId, g.Win, g.WinResult, nil, from)
	cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(gameId))
	if !ok {
		cacheData = &cache.CacheData{}
	}
	cacheData.IsEnd = true
	notifyGameEnd(gameId, g.Win, g.WinResult, cacheData, b, nil)
	handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
		GameId:      gameId,
		HandlerTime: time.Now().Unix(),
		Handler:     "end",
		Content:     b.GetCommitSGF(),
	}}, time.Second*5)
}

func makeRandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func SendUserOnline(gameId uint) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("SendUserOnline", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	blackOldEnter := false
	blackOldOnline := false
	whiteOldEnter := false
	whiteOldOnline := false
	lockKey := fmt.Sprintf("game-node:userOnlineTicker:%d", gameId)
	locker, err := redisUtils.NewLocker([]redis.Cmdable{redisUtils.Pool}, redisUtils.Options{
		KeyPrefix:   "lock",
		LockTimeout: 2 * time.Hour,
		WaitTimeout: 0,
		WaitRetry:   0,
	})
	if err != nil {
		return
	}
	lock, err := locker.Lock(lockKey)
	if err != nil {
		return
	}
	defer lock.Unlock()
	ticker := time.NewTicker(time.Second * 4)
	defer ticker.Stop()
	cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(gameId))
	if !ok {
		return
	}

	for {
		<-ticker.C
		if cacheData.IsEnd {
			return
		}
		if blackOldEnter != cacheData.BlackUserEnter || blackOldOnline != cacheData.BlackUserOnline {
			d := notifyStruct.UserStatusMsg{
				GameId:       gameId,
				UserId:       cacheData.BlackUserId,
				UserHash:     cacheData.BlackUserHash,
				EnterStatus:  cacheData.BlackUserEnter,
				OnlineStatus: cacheData.BlackUserOnline,
			}
			ws.GroupPublishApiChan <- &ws.GroupPublishData{
				GroupId: fmt.Sprintf("game:%d", gameId),
				Message: notifyStruct.MakeWsUserStatusMsg(&d),
				Ctx:     nil,
			}
			blackOldEnter = cacheData.BlackUserEnter
			blackOldOnline = cacheData.BlackUserOnline
		}
		if whiteOldEnter != cacheData.WhiteUserEnter || whiteOldOnline != cacheData.WhiteUserOnline {
			d := notifyStruct.UserStatusMsg{
				GameId:       gameId,
				UserId:       cacheData.WhiteUserId,
				UserHash:     cacheData.WhiteUserHash,
				EnterStatus:  cacheData.WhiteUserEnter,
				OnlineStatus: cacheData.WhiteUserOnline,
			}
			ws.GroupPublishApiChan <- &ws.GroupPublishData{
				GroupId: fmt.Sprintf("game:%d", gameId),
				Message: notifyStruct.MakeWsUserStatusMsg(&d),
				Ctx:     nil,
			}
			whiteOldEnter = cacheData.WhiteUserEnter
			whiteOldOnline = cacheData.WhiteUserOnline
		}
	}
}

func WithNotEnterTicker(gameId uint) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("WithNotEnterTicker", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(gameId))
	if !ok {
		return
	}
	if cacheData.NotEnterTime == 0 {
		return
	}
	diffTime := cacheData.NotEnterTime - time.Now().Unix()
	if diffTime <= 0 {
		return
	}
	timer := time.NewTicker(time.Second * time.Duration(diffTime))
	defer timer.Stop()
	select {
	case <-timer.C:
		if cacheData.BlackUserEnter && cacheData.WhiteUserEnter {
			return
		}
		b, ok := wq.WQDB.Get(wq.GameIdKey(gameId))
		if !ok {
			return
		}
		b.Paused = true
		cacheData.IsEnd = true
		var win int
		var winResult string
		if !cacheData.BlackUserEnter && !cacheData.WhiteUserEnter {
			// 双方都没进
			win = 4
			winResult = "Abstain"
		} else if !cacheData.BlackUserEnter && cacheData.WhiteUserEnter {
			// 黑方没进
			win = 2
			winResult = "W+L"
		} else if cacheData.BlackUserEnter && !cacheData.WhiteUserEnter {
			// 白方没进
			win = 1
			winResult = "B+L"
		}
		shutdownCh := GetShutdownChan(gameId)
		if shutdownCh != nil {
			select {
			case shutdownCh <- struct{}{}:
			case <-time.After(2 * time.Second):
				logger.Logger("WithNotEnterTicker", logger.ERROR, nil, fmt.Sprintf("gameId:%d, shutdownCh is closed", gameId))
			}
		}
		db := database.GetInstance()
		tx := db.Begin()
		e := model.UpdateGames(tx,
			&model.Game{Model: gorm.Model{ID: gameId}},
			&model.Game{IsEnd: true, EndTime: time.Now().Unix(), WinResult: winResult, Win: win, SGF: b.GetSGF(), Step: b.Step, BlackCaptured: b.CapturesBy[wq.BLACK], WhiteCaptured: b.CapturesBy[wq.WHITE], IsShow: true})
		if e != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		sendEnd(gameId, win, winResult, nil, "WithNotEnterTicker")
		notifyGameEnd(gameId, win, winResult, cacheData, b, nil)
		handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
			GameId:      gameId,
			HandlerTime: time.Now().Unix(),
			Handler:     "end",
			Content:     b.GetCommitSGF(),
		}}, time.Second*5)
		cacheData.IsEnd = true
		return
	}
}
