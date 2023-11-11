package handlerHistory

type HandlerHistory struct {
	GameId      uint   `json:"game_id"`
	UserId      string `json:"user_id"`
	UserHash    string `json:"user_hash"`
	Handler     string `json:"handler"`
	Duration    int64  `json:"duration"`
	HandlerTime int64  `json:"handler_time"`
	Content     string `json:"content"`
}
