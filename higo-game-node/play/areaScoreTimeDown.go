package play

import (
	"fmt"
	"higo-game-node/api/ws"
	"higo-game-node/cache"
	"higo-game-node/config"
	"higo-game-node/logger"
	"higo-game-node/notifyStruct"
	"higo-game-node/sharedMap"
	"higo-game-node/wq"
	"time"
)

var areaScoreTimeDownApplyUserMap = sharedMap.New[string]()

func GetAreaScoreTimeDownApplyUser(gameId uint) string {
	value, ok := areaScoreTimeDownApplyUserMap.Get(fmt.Sprintf("%d", gameId))
	if !ok {
		return ""
	}
	return value
}

func SetAreaScoreTimeDownApplyUser(gameId uint, applyUserId string) {
	areaScoreTimeDownApplyUserMap.Set(fmt.Sprintf("%d", gameId), applyUserId)
}

func DelAreaScoreTimeDownApplyUser(gameId uint) {
	areaScoreTimeDownApplyUserMap.Remove(fmt.Sprintf("%d", gameId))
}

var areaScoreTimeDownRejectChanMap = sharedMap.New[chan struct{}]()

func GetAreaScoreTimeDownRejectChan(gameId uint) chan struct{} {
	value, ok := areaScoreTimeDownRejectChanMap.Get(fmt.Sprintf("%d", gameId))
	if !ok {
		return nil
	}
	return value
}

func SetAreaScoreTimeDownRejectChan(gameId uint, ch chan struct{}) {
	areaScoreTimeDownRejectChanMap.Set(fmt.Sprintf("%d", gameId), ch)
}

func DeleteAreaScoreTimeDownRejectChan(gameId uint) {
	areaScoreTimeDownRejectChanMap.Remove(fmt.Sprintf("%d", gameId))
}

var areaScoreTimeDownAgreeChanMap = sharedMap.New[chan struct{}]()

func GetAreaScoreTimeDownAgreeChan(gameId uint) chan struct{} {
	value, ok := areaScoreTimeDownAgreeChanMap.Get(fmt.Sprintf("%d", gameId))
	if !ok {
		return nil
	}
	return value
}

func SetAreaScoreTimeDownAgreeChan(gameId uint, ch chan struct{}) {
	areaScoreTimeDownAgreeChanMap.Set(fmt.Sprintf("%d", gameId), ch)
}

func DeleteAreaScoreTimeDownAgreeChan(gameId uint) {
	areaScoreTimeDownAgreeChanMap.Remove(fmt.Sprintf("%d", gameId))
}

func AreaScoreTimeDown(gameId uint, applyUserId string, applyUserHash string, opponentUserId string, opponentHash string, agreeCh chan struct{}, rejectCh chan struct{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("AreaScoreTimeDown", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
			return
		}
		DeleteAreaScoreTimeDownRejectChan(gameId)
		DeleteAreaScoreTimeDownAgreeChan(gameId)
		DelAreaScoreTimeDownApplyUser(gameId)
		cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(gameId))
		if ok {
			cacheData.LastAreaScoreTime = time.Now().Unix()
		}
	}()
	totalTime := config.Conf.Rule.AreaScoreTime
	b, ok := wq.WQDB.Get(wq.GameIdKey(gameId))
	if ok {
		d := notifyStruct.AreaScoreTimeMsg{
			GameId:        gameId,
			ApplyUseId:    applyUserId,
			ApplyUserHash: applyUserHash,
			OpponentId:    opponentUserId,
			OpponentHash:  opponentHash,
		}
		d.LeftTime = totalTime
		d.Status = 1
		ws.GroupPublishApiChan <- &ws.GroupPublishData{
			GroupId: fmt.Sprintf("game:%d", gameId),
			Message: notifyStruct.MakeWsAreaScoreTimeMsg(&d),
		}
		b.Paused = true
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-rejectCh:
				d.LeftTime = totalTime
				d.Status = 2
				ws.GroupPublishApiChan <- &ws.GroupPublishData{
					GroupId: fmt.Sprintf("game:%d", gameId),
					Message: notifyStruct.MakeWsAreaScoreTimeMsg(&d),
					Ctx:     nil,
				}
				b.Paused = false
				turnCh := GetTurnChan(gameId)
				if turnCh != nil {
					func() {
						defer func() {
							if err := recover(); err != nil {
								logger.Logger("AreaScoreTimeDown", logger.ERROR, nil, "turnCh is closed")
							}
						}()
						turnCh <- b.Player
					}()
				}
				return
			case <-agreeCh:
				d.LeftTime = totalTime
				d.Status = 3
				ws.GroupPublishApiChan <- &ws.GroupPublishData{
					GroupId: fmt.Sprintf("game:%d", gameId),
					Message: notifyStruct.MakeWsAreaScoreTimeMsg(&d),
					Ctx:     nil,
				}
				return
			case <-ticker.C:
				totalTime--
				if totalTime == 0 {
					logger.Logger("AreaScoreTimeDown", logger.INFO, nil, fmt.Sprintf("total_time: %d", totalTime))
					d.LeftTime = 0
					d.Status = 4
					ws.GroupPublishApiChan <- &ws.GroupPublishData{
						GroupId: fmt.Sprintf("game:%d", gameId),
						Message: notifyStruct.MakeWsAreaScoreTimeMsg(&d),
					}
					b.Paused = false
					turnCh := GetTurnChan(gameId)
					if turnCh != nil {
						func() {
							defer func() {
								if err := recover(); err != nil {
									logger.Logger("AreaScoreTimeDown", logger.ERROR, nil, "rejectCh is closed")
								}
							}()
							turnCh <- b.Player
						}()
					}
					return
				} else {
					d.LeftTime = totalTime
					d.Status = 1
					ws.GroupPublishApiChan <- &ws.GroupPublishData{
						GroupId: fmt.Sprintf("game:%d", gameId),
						Message: notifyStruct.MakeWsAreaScoreTimeMsg(&d),
					}
				}
			}
		}
	}
}

func FakeUserAgreeAreaScore(request *InfoRequest) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("FakeUserAgreeAreaScore", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", request.GameId, err))
		}
	}()
	ticker := time.NewTicker(time.Second * time.Duration(makeRandInt(5, 15)))
	defer ticker.Stop()
	select {
	case <-ticker.C:
		cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(request.GameId))
		if !ok {
			return
		}
		if makeRandInt(0, 50)%5 != 0 {
			res := &EndResponse{}
			e := agreeAreaScoreService(request, cacheData, res)
			if e != nil {
				logger.Logger("FakeUserAgreeAreaScore", logger.ERROR, e, fmt.Sprintf("%s", e.Error()))
			}
		} else {
			e := rejectAreaScoreService(request, cacheData)
			if e != nil {
				logger.Logger("FakeUserRejectAreaScoreService", logger.ERROR, e, fmt.Sprintf("%s", e.Error()))
			}
		}
		return
	}
}
