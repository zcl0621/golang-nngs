package pub

import (
	"encoding/json"
	"fmt"
	"higo-game-bus/api/gameLb"
	"higo-game-bus/logger"
	"higo-game-bus/redisUtils"
	"time"
)

var initGameInfoChannel = "game:bus:init"

type InitPubData struct {
	GameId       uint  `json:"game_id"`
	CanStartTime int64 `json:"can_start_time"`
}

func PubInitData(data *InitPubData) error {
	d, _ := json.Marshal(data)
	e := redisUtils.LPush(initGameInfoChannel, d)
	return e
}

func StartInitData() {
	go func() {
		for {
			d, e := redisUtils.LRPop(initGameInfoChannel)
			if e != nil {
				continue
			}
			var data InitPubData
			e = json.Unmarshal(d, &data)
			if e != nil {
				continue
			}
			if data.CanStartTime <= time.Now().Unix() || time.Now().Unix()-data.CanStartTime <= 300 {
				e := gameLb.InitInfoApi(&gameLb.InitRequest{
					GameId: data.GameId,
				})
				if e != nil {
					logger.Logger("StartInitData", logger.ERROR, e, fmt.Sprintf("gameLb.InitInfoApi error: %v data: %v", e, data))
				}
			} else {
				e = PubInitData(&data)
				if e != nil {
					logger.Logger("StartInitData", logger.ERROR, e, fmt.Sprintf("pub init data error: %v data: %v", e, data))
				}
			}
		}
	}()
}
