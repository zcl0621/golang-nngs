package cache

import (
	"fmt"
	"gorm.io/gorm"
	"higo-game-node/model"
	"higo-game-node/redisUtils"
	"higo-game-node/sharedMap"
	"higo-game-node/syncLock"
	"time"
)

var CachaDataMap = sharedMap.New[*CacheData]()

type CacheData struct {
	TimeDoneIsReady   bool     `json:"time_done_is_ready"`
	IsEnd             bool     `json:"is_end"`
	BoardSize         int      `json:"board_size"`
	WinCaptured       int      `json:"win_captured"`     // 赢棋子数
	EnableMoveTime    int      `json:"enable_move_time"` // 是否开启第一步倒计时 0 否 1 是
	MaxStep           int      `json:"max_step"`         // 最大步数
	TerritoryStep     int      `json:"territory_step"`
	BlackReturnStone  float64  `json:"black_return_stone"` //黑返子数
	BlackUserId       string   `json:"black_user_id"`
	BlackUserHash     string   `json:"black_user_hash"`
	BlackUserEnter    bool     `json:"black_user_enter"`
	BlackUserOnline   bool     `json:"black_user_online"`
	BlackUserType     int      `json:"black_user_type"`
	BlackUserExtra    string   `json:"black_user_extra"`
	BlackClients      []string `json:"black_clients"`
	BlackClientsLock  *syncLock.ChanMutex
	WhiteUserId       string   `json:"white_user_id"`
	WhiteUserHash     string   `json:"white_user_hash"`
	WhiteUserEnter    bool     `json:"white_user_enter"`
	WhiteUserType     int      `json:"white_user_type"`
	WhiteUserExtra    string   `json:"white_user_extra"`
	WhiteUserOnline   bool     `json:"white_user_online"`
	WhiteClients      []string `json:"white_clients"`
	WhiteClientsLock  *syncLock.ChanMutex
	BusinessType      string `json:"business_type"` //业务类型
	CanStartTime      int64  `json:"can_start_time"`
	NotEnterTime      int64  `json:"not_enter_time"` // 一方未进入或两方都未进入 强制结束时间
	LastAreaScoreTime int64  `json:"last_area_score_time"`
	LastSummationTime int64  `json:"last_summation_time"`
	SummationCount    int    `json:"summation_count"`
	LastMoveTime      int64  `json:"last_move_time"`
	NowMoveTime       int    `json:"now_move_time"`
	TimeDownIsStarted bool   `json:"time_down_is_started"`
	Lock              *syncLock.ChanMutex
	IsMove            int
}

func GetKey(gameId uint) string {
	return fmt.Sprintf("cache-data-%d", gameId)
}
func InitSetCacheData(gameId uint, data *CacheData) {
	if _, ok := CachaDataMap.Get(GetKey(gameId)); ok {
		return
	}
	CachaDataMap.Set(GetKey(gameId), data)
	go func(data *CacheData) {
		diff := data.NotEnterTime - time.Now().Unix()
		if diff <= 8*3600 {
			diff = 8 * 3600
		}
		ticker := time.NewTicker(time.Second * time.Duration(diff))
		defer ticker.Stop()
		select {
		case <-ticker.C:
			data.Lock = nil
			data = nil
			CachaDataMap.Remove(GetKey(gameId))
			return
		}
	}(data)
}

func MakeCacheData(gameId uint) *CacheData {
	oneGame, _ := model.SelectOneGame(nil, &model.Game{Model: gorm.Model{ID: gameId}})
	if oneGame != nil {
		data := &CacheData{
			IsEnd:            oneGame.IsEnd,
			BoardSize:        oneGame.BoardSize,
			WinCaptured:      oneGame.WinCaptured,
			MaxStep:          oneGame.MaxStep,
			TerritoryStep:    oneGame.TerritoryStep,
			EnableMoveTime:   oneGame.EnableMoveTime,
			BlackReturnStone: oneGame.BlackReturnStone,
			BlackUserId:      oneGame.BlackUserId,
			BlackUserHash:    oneGame.BlackUserHash,
			BlackUserEnter:   oneGame.BlackUserEnter,
			BlackUserType:    oneGame.BlackUserType,
			BlackUserExtra:   oneGame.BlackUserExtra,
			BlackClientsLock: syncLock.NewChanMutex(),
			WhiteUserId:      oneGame.WhiteUserId,
			WhiteUserHash:    oneGame.WhiteUserHash,
			WhiteUserEnter:   oneGame.WhiteUserEnter,
			WhiteUserType:    oneGame.WhiteUserType,
			WhiteUserExtra:   oneGame.WhiteUserExtra,
			WhiteClientsLock: syncLock.NewChanMutex(),
			BusinessType:     oneGame.BusinessType,
			CanStartTime:     oneGame.CanStartTime,
			NotEnterTime:     oneGame.NotEnterTime,
			Lock:             syncLock.NewChanMutex(),
			IsMove:           0,
		}
		if oneGame.BlackUserEnter {
			data.BlackUserOnline = true
		}
		if oneGame.WhiteUserEnter {
			data.WhiteUserEnter = true
		}
		InitSetCacheData(gameId, data)
		return data
	}
	return nil
}

var GAMESGFKEY = "game:sgf:%d"

func SetGameSgf(gameId uint, sgf string) {
	redisUtils.Set(fmt.Sprintf(GAMESGFKEY, gameId), []byte(sgf), 3600)
}

func GetGameSgf(gameId uint) (sgf string) {
	data, err := redisUtils.Get(fmt.Sprintf(GAMESGFKEY, gameId))
	if err == nil {
		sgf = string(data)
	}
	return
}
