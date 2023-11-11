package steam

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"higo-game-bus/database"
	"higo-game-bus/logger"
	"higo-game-bus/model"
	"higo-game-bus/redisUtils"
	"time"
)

var steamGroup = "game:bus"

var notifyUserEnterChannels = "game-node:notify-user-enter"

type UserEnterNotify struct {
	GameId       uint   `json:"game_id"`
	UserId       string `json:"user_id"`
	BusinessType string `json:"business_type"`
}

var notifyGameBeginChannels = "game-node:notify-game-begin"

type GameBeginNotify struct {
	GameId       uint   `json:"game_id"`
	BusinessType string `json:"business_type"`
}

var notifyGameEndChannels = "game-node:notify-game-end"

type GameEndNotify struct {
	GameId       uint   `json:"game_id"`
	BusinessType string `json:"business_type"`
	Win          int    `json:"win"`
	WinResult    string `json:"win_result"`
}

func SubGameNotify() {
	go subGameBeginNotify()
	go subGameUserEnterNotify()
	go subGameEndNotify()
}

func subGameBeginNotify() {
	c := fmt.Sprintf("c:%d", time.Now().UnixNano())
	for {
		handlerGameBeginNotify(c)
	}
}

func handlerGameBeginNotify(c string) {
	res, mId, err := redisUtils.XReadGroup(notifyGameBeginChannels, steamGroup, c)
	if err != nil {
		logger.Logger("steam.subGameBeginNotify", logger.ERROR, err, fmt.Sprintf("steam消息消费失败 e %s", err))
		return
	}
	var notify GameBeginNotify
	err = json.Unmarshal(res, &notify)
	if err == nil {
		tx := database.GetInstance().Begin()
		defer func() {
			if e := recover(); e != nil {
				tx.Rollback()
				logger.Logger("steam.subGameBeginNotify", logger.ERROR, nil, fmt.Sprintf("steam消息消费失败 %v e %s", notify, e))
			}
		}()
		model.UpdateBattles(tx, &model.Battle{GameId: notify.GameId}, &model.Battle{StartedAt: time.Now().Unix(), IsShow: true})
		tx.Commit()
	} else {
		logger.Logger("steam.subGameBeginNotify", logger.ERROR, err, fmt.Sprintf("steam消息消费失败 %s", res))
	}
	_ = redisUtils.XAck(notifyGameBeginChannels, steamGroup, mId)
}

func subGameUserEnterNotify() {
	c := fmt.Sprintf("c:%d", time.Now().UnixNano())
	for {
		handlerGameUserEnterNotify(c)
	}
}

func handlerGameUserEnterNotify(c string) {
	_, mId, err := redisUtils.XReadGroup(notifyUserEnterChannels, steamGroup, c)
	if err != nil {
		logger.Logger("steam.subGameEndNotify", logger.ERROR, nil, fmt.Sprintf("steam消息消费失败 e %s", err))
		return
	}
	_ = redisUtils.XAck(notifyUserEnterChannels, steamGroup, mId)
}

func subGameEndNotify() {
	c := fmt.Sprintf("c:%d", time.Now().UnixNano())
	for {
		handlerGameEndNotify(c)
	}
}

func handlerGameEndNotify(c string) {
	res, mId, err := redisUtils.XReadGroup(notifyGameEndChannels, steamGroup, c)
	if err != nil {
		logger.Logger("steam.subGameEndNotify", logger.ERROR, nil, fmt.Sprintf("steam消息消费失败 e %s", err))
		return
	}
	var notify GameEndNotify
	err = json.Unmarshal(res, &notify)
	if err == nil {
		tx := database.GetInstance().Begin()
		defer func() {
			if e := recover(); e != nil {
				tx.Rollback()
				logger.Logger("steam.subGameEndNotify", logger.ERROR, nil, fmt.Sprintf("steam消息消费失败 %v e %s", notify, e))
			}
		}()
		switch notify.Win {
		case 3:
			model.UpdateBattles(tx,
				&model.Battle{GameId: notify.GameId},
				&model.Battle{EndedAt: time.Now().Unix(), UserWin: 3, WinResult: notify.WinResult, IsShow: true},
				func(db *gorm.DB) *gorm.DB {
					return db.Select("ended_at", "user_win", "win_result", "is_show")
				})
		case 4:
			model.UpdateBattles(tx,
				&model.Battle{GameId: notify.GameId},
				&model.Battle{EndedAt: time.Now().Unix(), UserWin: 3, WinResult: notify.WinResult, IsShow: false},
				func(db *gorm.DB) *gorm.DB {
					return db.Select("ended_at", "user_win", "win_result", "is_show")
				})
		default:
			allBattles := model.SelectBattles(tx, &model.Battle{GameId: notify.GameId})
			if len(allBattles) > 0 {
				for i := range allBattles {
					var userWin int
					if allBattles[i].UserSide == 1 && notify.Win == 1 {
						userWin = 1
					}
					if allBattles[i].UserSide == 2 && notify.Win == 2 {
						userWin = 1
					}
					if allBattles[i].UserSide == 1 && notify.Win == 2 {
						userWin = 2
					}
					if allBattles[i].UserSide == 2 && notify.Win == 1 {
						userWin = 2
					}
					model.UpdateBattles(tx, &model.Battle{Model: gorm.Model{ID: allBattles[i].ID}}, &model.Battle{UserWin: userWin, WinResult: notify.WinResult, EndedAt: time.Now().Unix(), IsShow: true})
				}
			}
		}
		tx.Commit()
	} else {
		logger.Logger("steam.subGameEndNotify", logger.ERROR, err, fmt.Sprintf("steam消息消费失败 %s", res))
	}
	_ = redisUtils.XAck(notifyGameEndChannels, steamGroup, mId)
}
