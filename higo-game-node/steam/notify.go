package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"higo-game-node/logger"
	"higo-game-node/redisUtils"
	"time"
)

func StartNotifyChan() {
	NotifyUserEnterChan = make(chan *UserEnterNotify, 100)
	NotifyGameBeginChan = make(chan *GameBeginNotify, 100)
	NotifyGameEndChan = make(chan *GameEndNotify, 100)
	go userEnter()
	go gameBegin()
	go gameEnd()
}

func userEnter() {
	for {
		notify := <-NotifyUserEnterChan
		go func(n *UserEnterNotify) {
			notifyUserEnter(notify)
		}(notify)
	}
}

func gameBegin() {
	for {
		notify := <-NotifyGameBeginChan
		go func(n *GameBeginNotify) {
			notifyGameBegin(n)
		}(notify)
	}
}

func gameEnd() {
	for {
		notify := <-NotifyGameEndChan
		go func(n *GameEndNotify) {
			notifyGameEnd(n)
		}(notify)
	}
}

var NotifyUserEnterChan chan *UserEnterNotify
var notifyUserEnterChannels = "game-node:notify-user-enter"

type UserEnterNotify struct {
	GameId       uint            `json:"game_id"`
	UserId       string          `json:"user_id"`
	BusinessType string          `json:"business_type"`
	Ctx          context.Context `json:"-"`
}

func notifyUserEnter(notify *UserEnterNotify) {
	d, _ := json.Marshal(*notify)
	if notify.Ctx == nil {
		if e := redisUtils.XAdd(notifyUserEnterChannels, d); e != nil {
			logger.Logger("StartNotifyChan.notifyUserEnter", logger.ERROR, e, fmt.Sprintf("steam消息生产失败 %v e %s", d, e.Error()))
		}
		return
	}
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	select {
	case <-notify.Ctx.Done():
		if e := redisUtils.XAdd(notifyUserEnterChannels, d); e != nil {
			logger.Logger("StartNotifyChan.notifyUserEnter", logger.ERROR, e, fmt.Sprintf("steam消息生产失败 %v e %s", d, e.Error()))
		}
		return
	case <-ticker.C:
		logger.Logger("StartNotifyChan.notifyUserEnter", logger.ERROR, nil, fmt.Sprintf("steam消息生产超时 %v", d))
		return
	}
}

var NotifyGameBeginChan chan *GameBeginNotify
var notifyGameBeginChannels = "game-node:notify-game-begin"

type GameBeginNotify struct {
	GameId       uint            `json:"game_id"`
	BusinessType string          `json:"business_type"`
	Ctx          context.Context `json:"-"`
}

func notifyGameBegin(notify *GameBeginNotify) {
	d, _ := json.Marshal(*notify)
	if notify.Ctx == nil {
		if e := redisUtils.XAdd(notifyGameBeginChannels, d); e != nil {
			logger.Logger("StartNotifyChan.notifyGameBegin", logger.ERROR, e, fmt.Sprintf("steam消息生产失败 %v e %s", d, e.Error()))
		}
		return
	}
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	select {
	case <-notify.Ctx.Done():
		if e := redisUtils.XAdd(notifyGameBeginChannels, d); e != nil {
			logger.Logger("StartNotifyChan.notifyGameBegin", logger.ERROR, e, fmt.Sprintf("steam消息生产失败 %v e %s", d, e.Error()))
		}
		return
	case <-ticker.C:
		logger.Logger("StartNotifyChan.notifyGameBegin", logger.ERROR, nil, fmt.Sprintf("steam消息生产超时 %v", d))
		return
	}
}

var NotifyGameEndChan chan *GameEndNotify
var notifyGameEndChannels = "game-node:notify-game-end"

type GameEndNotify struct {
	GameId            uint            `json:"game_id"`
	BusinessType      string          `json:"business_type"`
	Win               int             `json:"win"` // 0 无 1 黑 2 白 3 和
	WinResult         string          `json:"win_result"`
	BlackUserId       string          `json:"black_user_id"`
	BlackUserType     int             `json:"black_user_type"` //1 玩家 2 假人 3 AI
	BlackUserNickName string          `json:"black_user_nick_name"`
	BlackUserAvatar   string          `json:"black_user_avatar"`
	WhiteUserId       string          `json:"white_user_id"`
	WhiteUserType     int             `json:"white_user_type"` //1 玩家 2 假人 3 AI
	WhiteUserNickName string          `json:"white_user_nick_name"`
	WhiteUserAvatar   string          `json:"white_user_avatar"`
	WinCaptured       int             `json:"win_captured"` // 赢棋子数
	BCaptured         int             `json:"b_captured"`
	WCaptured         int             `json:"w_captured"`
	BScore            float64         `json:"b_score"`
	WScore            float64         `json:"w_score"`
	Step              int             `json:"step"`
	Ctx               context.Context `json:"-"`
}

func notifyGameEnd(notify *GameEndNotify) {
	d, _ := json.Marshal(*notify)
	if notify.Ctx == nil {
		if e := redisUtils.XAdd(notifyGameEndChannels, d); e != nil {
			logger.Logger("StartNotifyChan.notifyGameEnd", logger.ERROR, e, fmt.Sprintf("steam消息生产失败 %v e %s", d, e.Error()))
		}
		return
	}
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	select {
	case <-notify.Ctx.Done():
		if e := redisUtils.XAdd(notifyGameEndChannels, d); e != nil {
			logger.Logger("StartNotifyChan.notifyGameEnd", logger.ERROR, e, fmt.Sprintf("steam消息生产失败 %v e %s", d, e.Error()))
		}
		return
	case <-ticker.C:
		logger.Logger("StartNotifyChan.notifyGameEnd", logger.ERROR, nil, fmt.Sprintf("steam消息生产超时 %v", d))
		return
	}
}
