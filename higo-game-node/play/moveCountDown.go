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
	"higo-game-node/sharedMap"
	"higo-game-node/wq"
	"sync"
	"time"
)

var moveTurnChanMap = sharedMap.New[chan wq.Colour]()

func GetMoveTurnChan(gameId uint) chan wq.Colour {
	value, ok := moveTurnChanMap.Get(fmt.Sprintf("%d", gameId))
	if !ok {
		return nil
	}
	return value
}

func SetMoveTurnChan(gameId uint, turnCh chan wq.Colour) {
	moveTurnChanMap.Set(fmt.Sprintf("%d", gameId), turnCh)
}

func DeleteMoveTurnChan(gameId uint) {
	moveTurnChanMap.Remove(fmt.Sprintf("%d", gameId))
}

var moveTickEndChanMap = sharedMap.New[chan struct{}]()

func GetMoveTickEndChan(gameId uint) chan struct{} {
	value, ok := moveTickEndChanMap.Get(fmt.Sprintf("%d", gameId))
	if !ok {
		return nil
	}
	return value
}

func SetMoveTickerEndChan(gameId uint, turnCh chan struct{}) {
	moveTickEndChanMap.Set(fmt.Sprintf("%d", gameId), turnCh)
}

func DeleteTickerEndChan(gameId uint) {
	moveTickEndChanMap.Remove(fmt.Sprintf("%d", gameId))
}

func MoveTicker(gameId uint, timeData *gameTime, doneFunc func()) {
	once := sync.Once{}
	defer func() {
		if err := recover(); err != nil {
			once.Do(doneFunc)
			logger.Logger("MoveTicker", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
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
	var moveTurnCh chan wq.Colour
	moveTurnCh = make(chan wq.Colour)
	var moveTickerEndCh chan struct{}
	moveTickerEndCh = make(chan struct{})
	SetMoveTurnChan(gameId, moveTurnCh)
	SetMoveTickerEndChan(gameId, moveTickerEndCh)
	defer DeleteMoveTurnChan(gameId)
	defer DeleteTickerEndChan(gameId)

	if (timeData.BlackTime == 0 && timeData.BlackByoYomi == 0 && timeData.BlackByoYomiTime == 0) ||
		(timeData.WhiteTime == 0 && timeData.WhiteByoYomi == 0 && timeData.WhiteByoYomiTime == 0) {
		once.Do(doneFunc)
		return
	}

	count := 0
	b, ok := wq.WQDB.Get(wq.GameIdKey(gameId))
	if !ok {
		once.Do(doneFunc)
		return
	}
	ticker := time.NewTicker(time.Second * 1)
	totalTime := config.Conf.Rule.MoveTime
	var msg notifyStruct.MoveTimeMsg
	msg.GameId = gameId
	defer ticker.Stop()
	for {
		once.Do(doneFunc)
		select {
		case <-ticker.C:
			if b.Step == 0 {
				continue
			}
			if d.IsEnd {
				return
			}
			if b.Paused {
				break
			} else {
				totalTime--
				d.NowMoveTime = totalTime
				if totalTime <= 30 {
					count = 0
				} else {
					count++
				}
			}
			if totalTime > 0 {
				if count%10 == 0 {
					msg.Color = b.Player
					msg.LeftTime = totalTime
					ws.GroupPublishApiChan <- &ws.GroupPublishData{
						GroupId: fmt.Sprintf("game:%d", gameId),
						Message: notifyStruct.MakeWsMoveTimeMsg(&msg),
					}
				}
				ticker.Reset(time.Second * 1)
			} else {
				db := database.GetInstance()
				tx := db.Begin()
				b.Paused = true
				var win int
				var winResult string
				switch b.Player {
				case wq.BLACK:
					win = 2
					winResult = "W+T"
				case wq.WHITE:
					win = 1
					winResult = "B+T"
				}
				shutdownCh := GetShutdownChan(gameId)
				if shutdownCh != nil {
					select {
					case shutdownCh <- struct{}{}:
					case <-time.After(time.Second):
						logger.Logger("MoveTicker", logger.ERROR, nil, fmt.Sprintf("gameId:%d, shutdownCh is closed", gameId))
					}
				}
				e := model.UpdateGames(tx,
					&model.Game{Model: gorm.Model{ID: gameId}},
					&model.Game{IsEnd: true, EndTime: time.Now().Unix(), WinResult: winResult, Win: win, SGF: b.GetSGF(), Step: b.Step, BlackCaptured: b.CapturesBy[wq.BLACK], WhiteCaptured: b.CapturesBy[wq.WHITE], IsShow: true})
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
		case <-moveTickerEndCh:
			return
		case color := <-moveTurnCh:
			switch color {
			case wq.BLACK:
				bPlayer := GetPlayer(gameId, wq.BLACK)
				if bPlayer == nil {
					return
				}
				bLeftTime := bPlayer.mainTime + bPlayer.byoYomiPeriod*bPlayer.byoYomiTime
				if bLeftTime >= config.Conf.Rule.MoveTime {
					totalTime = config.Conf.Rule.MoveTime
				} else {
					totalTime = bLeftTime
				}
			case wq.WHITE:
				wPlayer := GetPlayer(gameId, wq.WHITE)
				if wPlayer == nil {
					return
				}
				wLeftTime := wPlayer.mainTime + wPlayer.byoYomiPeriod*wPlayer.byoYomiTime
				if wLeftTime >= config.Conf.Rule.MoveTime {
					totalTime = config.Conf.Rule.MoveTime
				} else {
					totalTime = wLeftTime
				}
			}
			count = 0
			msg.Color = b.Player
			msg.LeftTime = totalTime
			ws.GroupPublishApiChan <- &ws.GroupPublishData{
				GroupId: fmt.Sprintf("game:%d", gameId),
				Message: notifyStruct.MakeWsMoveTimeMsg(&msg),
			}
		}
	}
}
