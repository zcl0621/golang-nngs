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

var summationTimeDownApplyUserMap = sharedMap.New[string]()

func GetSummationTimeDownApplyUser(gameId uint) string {
	value, ok := summationTimeDownApplyUserMap.Get(fmt.Sprintf("%d", gameId))
	if !ok {
		return ""
	}
	return value
}

func SetSummationTimeDownApplyUser(gameId uint, applyUserId string) {
	summationTimeDownApplyUserMap.Set(fmt.Sprintf("%d", gameId), applyUserId)
}

func DelSummationTimeDownApplyUser(gameId uint) {
	summationTimeDownApplyUserMap.Remove(fmt.Sprintf("%d", gameId))
}

var summationTimeDownRejectChanMap = sharedMap.New[chan struct{}]()

func GetSummationTimeDownRejectChan(gameId uint) chan struct{} {
	value, ok := summationTimeDownRejectChanMap.Get(fmt.Sprintf("%d", gameId))
	if !ok {
		return nil
	}
	return value
}

func SetSummationTimeDownRejectChan(gameId uint, ch chan struct{}) {
	summationTimeDownRejectChanMap.Set(fmt.Sprintf("%d", gameId), ch)
}

func DeleteSummationTimeDownRejectChan(gameId uint) {
	summationTimeDownRejectChanMap.Remove(fmt.Sprintf("%d", gameId))
}

var summationTimeDownAgreeChanMap = sharedMap.New[chan struct{}]()

func GetSummationTimeDownAgreeChan(gameId uint) chan struct{} {
	value, ok := summationTimeDownAgreeChanMap.Get(fmt.Sprintf("%d", gameId))
	if !ok {
		return nil
	}
	return value
}

func SetSummationTimeDownAgreeChan(gameId uint, ch chan struct{}) {
	summationTimeDownAgreeChanMap.Set(fmt.Sprintf("%d", gameId), ch)
}

func DeleteSummationTimeDownAgreeChan(gameId uint) {
	summationTimeDownAgreeChanMap.Remove(fmt.Sprintf("%d", gameId))
}

func SummationTimeDown(gameId uint, applyUserId string, applyUserHash string, opponentUserId string, opponentHash string, agreeCh chan struct{}, rejectCh chan struct{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("SummationTimeDown", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", gameId, err))
			return
		}
		DeleteSummationTimeDownRejectChan(gameId)
		DeleteSummationTimeDownAgreeChan(gameId)
		DelSummationTimeDownApplyUser(gameId)
		cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(gameId))
		if ok {
			cacheData.LastSummationTime = time.Now().Unix()
		}
	}()
	totalTime := config.Conf.Rule.SummationTime
	b, ok := wq.WQDB.Get(wq.GameIdKey(gameId))
	if ok {
		d := notifyStruct.SummationTimeMsg{
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
			Message: notifyStruct.MakeWsSummationTimeMsg(&d),
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
					Message: notifyStruct.MakeWsSummationTimeMsg(&d),
					Ctx:     nil,
				}
				b.Paused = false
				turnCh := GetTurnChan(gameId)
				if turnCh != nil {
					func() {
						defer func() {
							if err := recover(); err != nil {
								logger.Logger("SummationTimeDown", logger.ERROR, nil, "turnCh is closed")
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
					Message: notifyStruct.MakeWsSummationTimeMsg(&d),
					Ctx:     nil,
				}
				return
			case <-ticker.C:
				totalTime--
				if totalTime == 0 {
					logger.Logger("SummationTimeDown", logger.INFO, nil, fmt.Sprintf("total_time: %d", totalTime))
					d.LeftTime = 0
					d.Status = 4
					ws.GroupPublishApiChan <- &ws.GroupPublishData{
						GroupId: fmt.Sprintf("game:%d", gameId),
						Message: notifyStruct.MakeWsSummationTimeMsg(&d),
					}
					b.Paused = false
					turnCh := GetTurnChan(gameId)
					if turnCh != nil {
						func() {
							defer func() {
								if err := recover(); err != nil {
									logger.Logger("SummationTimeDown", logger.ERROR, nil, "rejectCh is closed")
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
						Message: notifyStruct.MakeWsSummationTimeMsg(&d),
					}
				}
			}
		}
	}
}

func FakeUserAgreeSummation(request *InfoRequest) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("FakeUserAgreeSummation", logger.ERROR, nil, fmt.Sprintf("gameId:%d, err:%v", request.GameId, err))
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
			reply := &EndResponse{}
			e := agreeSummationService(request, cacheData, reply)
			if e != nil {
				logger.Logger("FakeUserAgreeSummationService", logger.ERROR, e, fmt.Sprintf("%s", e.Error()))
			}
		} else {
			e := rejectSummationService(request, cacheData)
			if e != nil {
				logger.Logger("FakeUserRejectSummationService", logger.ERROR, e, fmt.Sprintf("%s", e.Error()))
			}
		}
		return
	}
}
