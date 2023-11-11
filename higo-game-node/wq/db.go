package wq

import (
	"fmt"
	"gorm.io/gorm"
	"higo-game-node/cache"
	"higo-game-node/model"
	"higo-game-node/sharedMap"
	"time"
)

var WQDB = sharedMap.New[*Board]()

func GameIdKey(id uint) string {
	return fmt.Sprintf("board-%d", id)
}

func MakeWq(id uint, endTime int64) (*Board, error) {
	g, ok := WQDB.Get(GameIdKey(id))
	if !ok {
		var thisSgf string
		redisSGF := cache.GetGameSgf(id)
		if redisSGF != "" {
			thisSgf = redisSGF
		} else {
			oneGame, _ := model.SelectOneGame(nil, &model.Game{Model: gorm.Model{ID: id}})
			if oneGame.SGF == "" {
				oneGame.SGF = fmt.Sprintf("(;SZ[%d]KM[7.5])", oneGame.BoardSize)
			}
			thisSgf = oneGame.SGF
		}

		return SetWq(id, thisSgf, endTime)
	}
	return g, nil
}

func SetWq(gameID uint, sgf string, endTime int64) (*Board, error) {
	newGame, _, err := Load(sgf)
	if err != nil {
		return nil, err
	}
	b := newGame.Board()
	WQDB.Set(GameIdKey(gameID), b)
	go func(gameId uint, board *Board) {
		diff := endTime - time.Now().Unix()
		if diff <= 8*3600 {
			diff = 8 * 3600
		}
		ticker := time.NewTimer(time.Second * time.Duration(diff))
		defer ticker.Stop()
		select {
		case <-ticker.C:
			board = nil
			WQDB.Remove(GameIdKey(gameId))
			return
		}
	}(gameID, b)
	return b, nil
}
