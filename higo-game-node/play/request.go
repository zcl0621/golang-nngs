package play

import "higo-game-node/wq"

type MoveRequest struct {
	BaseRequest
	X int       `json:"x"`
	Y int       `json:"y"`
	C wq.Colour `json:"c"`
}

type BaseRequest struct {
	GameId   uint   `json:"game_id" binding:"required"`
	UserId   string `json:"user_id" binding:"required"`
	UserHash string `json:"user_hash" binding:"required"`
}

type InfoRequest struct {
	BaseRequest
}

type EnterRequest struct {
	BaseRequest
}

type ResignRequest struct {
	BaseRequest
}

type PassRequest struct {
	BaseRequest
	C wq.Colour `json:"c"`
}

type InitRequest struct {
	GameId uint `json:"game_id"`
}

type OwnershipRequest struct {
	GameId uint   `json:"game_id"`
	SGF    string `json:"sgf"`
}

type CallEndRequest struct {
	GameId    uint   `json:"game_id" binding:"required"`
	SGF       string `json:"sgf"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Win       int    `json:"win" binding:"required"`        // 1 黑 2 白 3 和
	WinResult string `json:"win_result" binding:"required"` // 输赢原因 B+R 黑胜白认输 B+T 黑胜白超时 W+R 白胜黑认输 W+T 白胜黑超时 B+C1 黑胜吃1子 W+C1 白胜吃一子 B+L 黑胜白未参加 W+L 白胜黑未参加 Draw 和
}

type ForceReloadSgfRequest struct {
	GameId uint   `json:"game_id" binding:"required"`
	SGF    string `json:"sgf" binding:"required"`
}
