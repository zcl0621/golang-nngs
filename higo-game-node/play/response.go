package play

import (
	"higo-game-node/api/gameBus"
	"higo-game-node/wq"
)

type InfoResponse struct {
	BoardSize           int     `json:"board_size"`     //9 13 19
	NotEnterTime        int64   `json:"not_enter_time"` // 一方未进入或两方都未进入 强制结束时间
	IsStart             bool    `json:"is_start"`
	StartTime           int64   `json:"start_time"`
	IsEnd               bool    `json:"is_end"`
	EndTime             int64   `json:"end_time"`
	Win                 int     `json:"win"`        // 0 无 1 黑 2 白 3 和 4 弃权
	WinResult           string  `json:"win_result"` // 输赢原因 B+2.5 黑胜2.5目 W+2.5 白胜2.5目 B+R 黑胜白认输 B+T 黑胜白超时 W+R 白胜黑认输 W+T 白胜黑超时 B+C1 黑胜吃1子 W+C1 白胜吃一子 B+L 黑胜白未参加 W+L  白胜黑未参加 B+A 黑胜仲裁 W+A 白胜仲裁 Draw 和 Abstain 弃权
	WinCapture          int     `json:"win_capture"`
	KM                  float64 `json:"km"`               //贴目
	MoveTime            int     `json:"move_time"`        //单步落子时长
	SummationCount      int     `json:"summation_count"`  //申请和棋数量
	EnableMoveTime      int     `json:"enable_move_time"` // 是否开启第一步倒计时 0 否 1 是
	MaxStep             int     `json:"max_step"`
	TerritoryStep       int     `json:"territory_step"`
	BlackReturnStone    float64 `json:"black_return_stone"`
	Step                int     `json:"step"`
	BlackTime           int     `json:"black_time"`
	WhiteTime           int     `json:"white_time"`
	BlackByoYomi        int     `json:"black_byo_yomi"`
	WhiteByoYomi        int     `json:"white_byo_yomi"`
	BlackByoYomiTime    int     `json:"black_byo_yomi_time"`
	WhiteByoYomiTime    int     `json:"white_byo_yomi_time"`
	BlackUserId         string  `json:"black_user_id"`
	BlackUserName       string  `json:"black_user_name"`
	BlackUserAccount    string  `json:"black_user_account"`
	BlackUserNickName   string  `json:"black_user_nick_name"`
	BlackUserActualName string  `json:"black_user_actual_name"`
	BlackUserAvatar     string  `json:"black_user_avatar"`
	BlackUserLevel      string  `json:"black_user_level"`
	BlackUserHash       string  `json:"black_user_hash"`
	BlackUserType       int     `json:"black_user_type"`
	BlackUserEnter      bool    `json:"black_enter"`
	BlackUserOnline     bool    `json:"black_user_online"`
	WhiteUserId         string  `json:"white_user_id"`
	WhiteUserName       string  `json:"white_user_name"`
	WhiteUserAccount    string  `json:"white_user_account"`
	WhiteUserNickName   string  `json:"white_user_nick_name"`
	WhiteUserActualName string  `json:"white_user_actual_name"`
	WhiteUserAvatar     string  `json:"white_user_avatar"`
	WhiteUserLevel      string  `json:"white_user_level"`
	WhiteUserHash       string  `json:"white_user_hash"`
	WhiteUserType       int     `json:"white_user_type"`
	WhiteUserEnter      bool    `json:"white_enter"`
	WhiteUserOnline     bool    `json:"white_user_online"`
	BusinessType        string  `json:"business_type"` //业务类型

	BlackCaptured int     `json:"black_captured"`
	WhiteCaptured int     `json:"white_captured"`
	BlackScore    float64 `json:"black_score"`
	WhiteScore    float64 `json:"white_score"`

	NowMoveTime         int       `json:"now_move_time"`
	NowBlackTime        int       `json:"now_black_time"`
	NowWhiteTime        int       `json:"now_white_time"`
	NowBlackByoYomi     int       `json:"now_black_byo_yomi"`
	NowWhiteByoYomi     int       `json:"now_white_byo_yomi"`
	NowBlackByoYomiTime int       `json:"now_black_byo_yomi_time"`
	NowWhiteByoYomiTime int       `json:"now_white_byo_yomi_time"`
	Turn                wq.Colour `json:"turn"`
}

type MoveResponse struct {
	NextColor wq.Colour `json:"next_color"` // 下一步颜色 0 任意 1 黑 2 白
}

type PassResponse struct {
	NextColor wq.Colour `json:"next_color"` // 下一步颜色 0 任意 1 黑 2 白
}

type CanPlayResponse struct {
	NextColor wq.Colour `json:"next_color"`
}

type EndResponse struct {
	Win       int    `json:"win"`        // 0 无 1 黑 2 白
	WinResult string `json:"win_result"` // 输赢原因
}

type SGFResponse struct {
	SGF string `json:"sgf"`
}

type OwnerShipResponse struct {
	Ownership    []wq.OwnerShip            `json:"ownership"`
	AnalysisData gameBus.AnalysisScoreData `json:"analysis_data"`
}

type InnerAreaScoreResponse struct {
	BScore           float64 `json:"b_score"`
	WScore           float64 `json:"w_score"`
	EndScore         float64 `json:"end_score"`
	ControversyCount int     `json:"controversy_count"`
}
