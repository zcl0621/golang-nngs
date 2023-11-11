package public

type battleResponse struct {
	StartedAt              int64  `json:"started_at"`
	EndedAt                int64  `json:"ended_at"`
	GameId                 uint   `json:"game_id"`
	BoardSize              int    `json:"board_size"`
	Type                   string `json:"type"` // territory  captured
	UserName               string `json:"user_name"`
	UserAccount            string `json:"user_account"`
	UserNickName           string `json:"user_nick_name"`
	UserActualName         string `json:"user_actual_name"`
	UserAvatar             string `json:"user_avatar"`
	UserSide               int    `json:"user_side"`  // 1 黑 2 白
	UserWin                int    `json:"user_win"`   // 1 胜 2 负 3 和 4 弃权
	WinResult              string `json:"win_result"` // 输赢原因 B+2.5 黑胜2.5目 W+2.5 白胜2.5目 B+R 黑胜白认输 B+T 黑胜白超时 W+R 白胜黑认输 W+T 白胜黑超时 B+C1 黑胜吃1子 W+C1 白胜吃一子 B+L 黑胜白未参加 W+L  白胜黑未参加 B+A 黑胜仲裁 W+A 白胜仲裁 Draw 和 Abstain 弃权
	OpponentUserName       string `json:"opponent_user_name"`
	OpponentUserAccount    string `json:"opponent_user_account"`
	OpponentUserNickName   string `json:"opponent_user_nick_name"`
	OpponentUserActualName string `json:"opponent_user_actual_name"`
	OpponentUserAvatar     string `json:"opponent_user_avatar"`
	OpponentUserSide       int    `json:"opponent_user_side"` // 1 黑 2 白
	BusinessType           string `json:"business_type"`      //业务类型
}

type GameResponse struct {
	Id                  uint    `json:"id"`
	CanStartTime        int64   `json:"can_start_time"`
	NotEnterTime        int64   `json:"not_enter_time"` // 一方未进入或两方都未进入 强制结束时间
	IsStart             bool    `json:"is_start"`
	StartTime           int64   `json:"start_time"`
	IsEnd               bool    `json:"is_end"`
	EndTime             int64   `json:"end_time"`
	WinCaptured         int     `json:"win_captured"` // 赢棋子数
	KM                  float64 `json:"km"`           //贴目
	MaxStep             int     `json:"max_step"`     // 最大步数
	Win                 int     `json:"win"`          // 0 无 1 黑 2 白 3 和 4 弃权
	WinResult           string  `json:"win_result"`   // 输赢原因 B+2.5 黑胜2.5目 W+2.5 白胜2.5目 B+R 黑胜白认输 B+T 黑胜白超时 W+R 白胜黑认输 W+T 白胜黑超时 B+C1 黑胜吃1子 W+C1 白胜吃一子 B+L 黑胜白未参加 W+L  白胜黑未参加 B+A 黑胜仲裁 W+A 白胜仲裁 Draw 和 Abstain 弃权
	Step                int     `json:"step"`
	BlackUserId         string  `json:"black_user_id"`
	BlackUserName       string  `json:"black_user_name"`
	BlackUserAccount    string  `json:"black_user_account"`
	BlackUserNickName   string  `json:"black_user_nick_name"`
	BlackUserActualName string  `json:"black_user_actual_name"`
	BlackUserAvatar     string  `json:"black_user_avatar"`
	BlackUserLevel      string  `json:"black_user_level"`
	BlackUserEnter      bool    `json:"black_enter"`
	BlackUserType       int     `json:"black_user_type"` // 1 玩家 2 假人 3 AI
	WhiteUserId         string  `json:"white_user_id"`
	WhiteUserName       string  `json:"white_user_name"`
	WhiteUserAccount    string  `json:"white_user_account"`
	WhiteUserNickName   string  `json:"white_user_nick_name"`
	WhiteUserActualName string  `json:"white_user_actual_name"`
	WhiteUserAvatar     string  `json:"white_user_avatar"`
	WhiteUserLevel      string  `json:"white_user_level"`
	WhiteUserEnter      bool    `json:"white_user_enter"`
	WhiteUserType       int     `json:"white_user_type"` // 1 玩家 2 假人 3 AI
	BusinessType        string  `json:"business_type"`   //业务类型

	BlackTime            int `json:"black_time"`               //	黑方初始时间
	WhiteTime            int `json:"white_time"`               // 白方初始时间
	BlackByoYomi         int `json:"black_byo_yomi"`           // 黑方初始读秒
	WhiteByoYomi         int `json:"white_byo_yomi"`           // 白方初始读秒
	BlackByoYomiTime     int `json:"black_byo_yomi_time"`      // 黑方初始读秒时间
	WhiteByoYomiTime     int `json:"white_byo_yomi_time"`      // 白方初始读秒时间
	LeftBlackTime        int `json:"left_black_time"`          // 黑方剩余时间
	LeftWhiteTime        int `json:"left_white_time"`          // 白方剩余时间
	LeftBlackByoYomiTime int `json:"left_black_byo_yomi_time"` // 黑方剩余读秒时间
	LeftWhiteByoYomiTime int `json:"left_white_byo_yomi_time"` // 白方剩余读秒时间
	LeftWhiteByoYomi     int `json:"left_white_byo_yomi"`      // 白方剩余读秒次数
	LeftBlackByoYomi     int `json:"left_black_byo_yomi"`      // 黑方剩余读秒次数

	BoardSize int `json:"board_size"` //棋盘尺寸
	ViewCount int `json:"view_count"` //观看人数
}

type SGFResponse struct {
	SGF string `json:"sgf"`
}

type businessTypeView struct {
	Time string         `json:"time"`
	Data []businessData `json:"data"`
}

type businessData struct {
	BusinessType string `json:"business_type"`
	Count        int    `json:"count"`
}
