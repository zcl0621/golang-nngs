package play

import (
	"fmt"
	"gorm.io/gorm"
	"higo-game-node/cache"
	"higo-game-node/model"
	"higo-game-node/redisUtils"
	"higo-game-node/utils"
	"higo-game-node/wq"
	"strconv"
	"strings"
	"sync"
	"time"
)

func AddPlayingGame(gameId uint) {
	hostName := utils.GetHostName()
	key := fmt.Sprintf("game:node:%s:%d:playing", hostName, gameId)
	redisUtils.Set(key, []byte(fmt.Sprintf("%d", gameId)), 4*60*60)
}

func DeletePlayingGame(gameId uint) {
	hostName := utils.GetHostName()
	key := fmt.Sprintf("game:node:%s:%d:playing", hostName, gameId)
	redisUtils.Del(key)
}

func GetAllPlayingGame() []uint {
	hostName := utils.GetHostName()
	key := fmt.Sprintf("game:node:%s:*:playing", hostName)
	keys, _ := redisUtils.Keys(key)
	var gameIds []uint
	for i := range keys {
		val, _ := redisUtils.Get(keys[i])
		gameId, _ := strconv.Atoi(string(val))
		gameIds = append(gameIds, uint(gameId))
	}
	return gameIds
}

func SetGamePlayerTimeToRedis(gameId uint, color wq.Colour, mainTime int, byoYomiTime int, byoYomiPeriod int) {
	key := fmt.Sprintf("game:%d:player:%d", gameId, color)
	redisUtils.Set(key, []byte(fmt.Sprintf("%d:%d:%d", mainTime, byoYomiTime, byoYomiPeriod)), 60)
}

func GetGamePlayerTimeFromRedis(gameId uint, color wq.Colour) (mainTime int, byoYomiTime int, byoYomiPeriod int, has bool) {
	key := fmt.Sprintf("game:%d:player:%d", gameId, color)
	val, err := redisUtils.Get(key)
	if err != nil {
		return
	}
	str := string(val)
	arr := strings.Split(str, ":")
	if len(arr) != 3 {
		return
	}
	mainTime, _ = strconv.Atoi(arr[0])
	byoYomiTime, _ = strconv.Atoi(arr[1])
	byoYomiPeriod, _ = strconv.Atoi(arr[2])
	has = true
	return
}

func ReloadGame() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		<-ticker.C
		break
	}
	gameIds := GetAllPlayingGame()
	for i := range gameIds {
		oneModel, _ := model.SelectOneGame(nil, &model.Game{Model: gorm.Model{ID: gameIds[i]}})
		if oneModel == nil {
			continue
		}
		if oneModel.IsEnd {
			DeletePlayingGame(gameIds[i])
			continue
		}
		e := initService(&InitRequest{
			GameId: gameIds[i],
		})
		if e != nil {
			continue
		}
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
		bMainTime, byoYomiTime, byoYomiPeriod, bhas := GetGamePlayerTimeFromRedis(gameIds[i], wq.BLACK)
		wMainTime, wbyoYomiTime, wbyoYomiPeriod, whas := GetGamePlayerTimeFromRedis(gameIds[i], wq.WHITE)
		if bhas && whas {
			timeData.BlackTime = bMainTime
			timeData.WhiteTime = wMainTime
			timeData.BlackByoYomi = byoYomiPeriod
			timeData.WhiteByoYomi = wbyoYomiPeriod
			timeData.BlackByoYomiTime = byoYomiTime
			timeData.WhiteByoYomiTime = wbyoYomiTime
		}
		wg := sync.WaitGroup{}
		wg.Add(2)
		go StartTimeDown(
			gameIds[i],
			timeData,
			func() {
				wg.Done()
			},
		)
		go MoveTicker(gameIds[i], timeData, func() {
			wg.Done()
		})
		wg.Wait()
		cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(gameIds[i]))
		if ok {
			cacheData.TimeDownIsStarted = true
		}
	}
}
