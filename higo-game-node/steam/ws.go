package steam

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"higo-game-node/api/ws"
	"higo-game-node/cache"
	"higo-game-node/database"
	"higo-game-node/logger"
	"higo-game-node/model"
	"higo-game-node/notifyStruct"
	"higo-game-node/redisUtils"
	"higo-game-node/utils"
	"regexp"
	"strings"
	"time"
)

var redisClientStatusChannels = "ws-gateway-channels:client-status"

type wsClientStatusMessage struct {
	GroupId  string `json:"group_id"`
	ClientId string `json:"client_id"`
	UserId   string `json:"user_id"`  //棋盘编号
	Type     string `json:"type"`     // 上线 enter 下线 leave 监听组 group_enter 取消监听组 group_leave
	Platform string `json:"platform"` // board app
}

func SubWsConnect() {
	_, err := redisUtils.Subscribe(redisClientStatusChannels, func(message *redis.Message, err error) {
		if err != nil {
			logger.Logger("steam.SubWsConnect", logger.ERROR, err, fmt.Sprintf("can't subscribe message value:%v", message))
			return
		}
		var data wsClientStatusMessage
		err = json.Unmarshal([]byte(message.Payload), &data)
		if err != nil {
			logger.Logger("steam.SubWsConnect", logger.ERROR, err, fmt.Sprintf("can't Unmarshal message value:%v", message))
			return
		}
		if data.Platform == "app" {
			match, err := regexp.MatchString(`^game`, data.GroupId)
			if err == nil && match {
				ids := strings.Split(data.GroupId, ":")
				if len(ids) == 2 {
					gameId := ids[1]
					hostname := utils.GetHostName()
					if hostname == utils.FindPodHost(gameId) {
						cacheData, ok := cache.CachaDataMap.Get(cache.GetKey(uint(utils.StringToInt(gameId))))
						if !ok {
							return
						}
						switch data.Type {
						case "group_enter":
							switch fmt.Sprintf("%s:1", data.UserId) {
							case cacheData.BlackUserHash:
								ok := cacheData.BlackClientsLock.TryLockWithTimeout(time.Second * 5)
								if ok {
									cacheData.BlackClients = append(cacheData.BlackClients, data.ClientId)
									cacheData.BlackUserOnline = true
									sendUserStatus(uint(utils.StringToInt(gameId)), cacheData)
									cacheData.BlackClientsLock.Unlock()
								}
							case cacheData.WhiteUserHash:
								if ok := cacheData.WhiteClientsLock.TryLockWithTimeout(time.Second * 5); ok {
									cacheData.WhiteClients = append(cacheData.WhiteClients, data.ClientId)
									cacheData.WhiteUserOnline = true
									sendUserStatus(uint(utils.StringToInt(gameId)), cacheData)
									cacheData.WhiteClientsLock.Unlock()
								}
							default:
								db := database.GetInstance()
								tx := db.Begin()
								e := model.ScanGames(tx, &model.Game{Model: gorm.Model{ID: uint(utils.StringToInt(gameId))}}, func(db *gorm.DB) *gorm.DB {
									return db.Update("view_count", gorm.Expr("view_count + ?", 1))
								})
								if e != nil {
									tx.Rollback()
								} else {
									tx.Commit()
								}
							}
						case "group_leave":
							switch fmt.Sprintf("%s:1", data.UserId) {
							case cacheData.BlackUserHash:
								if ok := cacheData.BlackClientsLock.TryLockWithTimeout(time.Second * 5); ok {
									for i, v := range cacheData.BlackClients {
										if v == data.ClientId {
											cacheData.BlackClients = append(cacheData.BlackClients[:i], cacheData.BlackClients[i+1:]...)
											break
										}
									}
									if len(cacheData.BlackClients) == 0 {
										cacheData.BlackUserOnline = false
										sendUserStatus(uint(utils.StringToInt(gameId)), cacheData)
									}
									cacheData.BlackClientsLock.Unlock()
								}
							case cacheData.WhiteUserHash:
								if ok := cacheData.WhiteClientsLock.TryLockWithTimeout(time.Second * 5); ok {
									for i, v := range cacheData.WhiteClients {
										if v == data.ClientId {
											cacheData.WhiteClients = append(cacheData.WhiteClients[:i], cacheData.WhiteClients[i+1:]...)
											break
										}
									}
									if len(cacheData.WhiteClients) == 0 {
										cacheData.WhiteUserOnline = false
										sendUserStatus(uint(utils.StringToInt(gameId)), cacheData)
									}
									cacheData.WhiteClientsLock.Unlock()
								}
							}
						}
					}
				}
			}
		}
	})
	if err != nil {
		panic(err)
	}
}

func sendUserStatus(gameId uint, cacheData *cache.CacheData) {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	select {
	case <-ticker.C:
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
		d = notifyStruct.UserStatusMsg{
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
		return
	}

}
