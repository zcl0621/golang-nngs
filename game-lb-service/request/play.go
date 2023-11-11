package request

type MoveRequest struct {
	BaseRequest
	X int `json:"x"`
	Y int `json:"y"`
	C int `json:"c"`
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
	C int `json:"c"`
}

type InitRequest struct {
	GameId uint `json:"game_id"`
}

type OwnershipRequest struct {
	GameId uint   `json:"game_id"`
	SGF    string `json:"sgf"`
}

type ForceReloadSgfRequest struct {
	GameId uint   `json:"game_id" binding:"required"`
	SGF    string `json:"sgf" binding:"required"`
}
