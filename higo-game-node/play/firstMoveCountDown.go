package play

import (
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
	"higo-game-node/wq"
	"sync"
	"time"
)

func StartFirstMoveCountDown(gameId uint, doneFunc func()) {
	once := sync.Once{}
	defer func() {
		if err := recover(); err != nil {
			once.Do(doneFunc)
			logger.Logger("StartFirstMoveCountDown", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
		}
	}()
	d, ok := cache.CachaDataMap.Get(cache.GetKey(gameId))
	if !ok {
		once.Do(doneFunc)
		return
	}
	if d.EnableMoveTime == 0 {
		once.Do(doneFunc)
		return
	}
	b, ok := wq.WQDB.Get(wq.GameIdKey(gameId))
	if !ok {
		once.Do(doneFunc)
		return
	}
	if b.Step >= 1 {
		once.Do(doneFunc)
		return
	}
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	totalTime := config.Conf.Rule.FirstMoveTime
	var msg notifyStruct.MoveTimeMsg
	for {
		once.Do(doneFunc)
		select {
		case <-ticker.C:
			if b.Step >= 1 {
				return
			}
			if d.IsEnd {
				return
			}
			if totalTime >= 0 {
				msg.Color = wq.BLACK
				msg.LeftTime = totalTime
				ws.GroupPublishApiChan <- &ws.GroupPublishData{
					GroupId: fmt.Sprintf("game:%d", gameId),
					Message: notifyStruct.MakeWsMoveTimeMsg(&msg),
				}
			}
			if totalTime <= 0 {
				db := database.GetInstance()
				tx := db.Begin()
				b.Paused = true
				win := 4
				winResult := "Abstain"
				shutdownCh := GetShutdownChan(gameId)
				if shutdownCh != nil {
					select {
					case shutdownCh <- struct{}{}:
					case <-time.After(time.Second):
						logger.Logger("StartFirstMoveCountDown", logger.ERROR, nil, fmt.Sprintf("gameId:%d, shutdownCh is closed", gameId))
					}
				}
				e := model.UpdateGames(tx,
					&model.Game{Model: gorm.Model{ID: gameId}},
					&model.Game{
						IsEnd:         true,
						EndTime:       time.Now().Unix(),
						WinResult:     winResult,
						Win:           win,
						SGF:           b.GetSGF(),
						Step:          b.Step,
						BlackCaptured: b.CapturesBy[wq.BLACK],
						WhiteCaptured: b.CapturesBy[wq.WHITE],
						IsShow:        false},
					func(db *gorm.DB) *gorm.DB {
						return db.Select("is_end", "end_time", "win_result", "win", "sgf", "step", "black_captured", "white_captured", "is_show")
					})
				if e != nil {
					b.Paused = false
					tx.Rollback()
					return
				}
				sendEnd(gameId, win, winResult, nil, "MoveTicker")
				notifyGameEnd(gameId, win, winResult, d, b, nil)
				handlerHistory.Queue.ScheduleMessage(&handlerHistory.Message{Content: &handlerHistory.HandlerHistory{
					GameId:      gameId,
					HandlerTime: time.Now().Unix(),
					Handler:     "end",
					Content:     b.GetCommitSGF(),
				}}, time.Second*5)
				d.IsEnd = true
				tx.Commit()
				return
			}
			totalTime--
		}
	}
}
