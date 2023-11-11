package inner

type createRequest struct {
	BoardSize           int     `json:"board_size" binding:"required"` //9 13 19
	SGF                 string  `json:"sgf" binding:"required"`
	CanStartTime        int64   `json:"can_start_time"`
	NotEnterTime        int64   `json:"not_enter_time"`   // 一方未进入或两方都未进入 强制结束时间
	EnableMoveTime      int     `json:"enable_move_time"` // 是否开启第一步倒计时 0 否 1 是
	IsStart             bool    `json:"is_start"`
	StartTime           int64   `json:"start_time"`
	IsEnd               bool    `json:"is_end"`
	EndTime             int64   `json:"end_time"`
	WinCaptured         int     `json:"win_captured"` // 赢棋子数
	MaxStep             int     `json:"max_step"`     // 最大步数
	TerritoryStep       int     `json:"territory_step"`
	BlackReturnStone    float64 `json:"black_return_stone"`
	Win                 int     `json:"win"`        // 0 无 1 黑 2 白 3 和 4 双方弃权
	WinResult           string  `json:"win_result"` // 输赢原因 B+2.5 黑胜2.5目 W+2.5 白胜2.5目 B+R 黑胜白认输 B+T 黑胜白超时 W+R 白胜黑认输 W+T 白胜黑超时 B+C1 黑胜吃1子 W+C1 白胜吃一子 B+L 黑胜白未参加 W+L  白胜黑未参加 B+A 黑胜仲裁 W+A 白胜仲裁 Draw 和 Abstain 弃权
	Step                int     `json:"step"`
	BlackTime           int     `json:"black_time"`
	WhiteTime           int     `json:"white_time"`
	BlackByoYomi        int     `json:"black_byo_yomi"`
	WhiteByoYomi        int     `json:"white_byo_yomi"`
	BlackByoYomiTime    int     `json:"black_byo_yomi_time"`
	WhiteByoYomiTime    int     `json:"white_byo_yomi_time"`
	BlackUserId         string  `json:"black_user_id" binding:"required"`
	BlackUserName       string  `json:"black_user_name"`
	BlackUserAccount    string  `json:"black_user_account"`
	BlackUserNickName   string  `json:"black_user_nick_name"`
	BlackUserActualName string  `json:"black_user_actual_name"`
	BlackUserAvatar     string  `json:"black_user_avatar"`
	BlackUserLevel      string  `json:"black_user_level"`
	BlackUserType       int     `json:"black_user_type" binding:"required"` // 1 玩家 2 假人 3 AI
	BlackUserExtra      string  `json:"black_user_extra"`                   // 附加信息
	BlackUserEnter      bool    `json:"black_enter"`
	WhiteUserId         string  `json:"white_user_id" binding:"required"`
	WhiteUserName       string  `json:"white_user_name"`
	WhiteUserAccount    string  `json:"white_user_account"`
	WhiteUserNickName   string  `json:"white_user_nick_name"`
	WhiteUserActualName string  `json:"white_user_actual_name"`
	WhiteUserAvatar     string  `json:"white_user_avatar"`
	WhiteUserLevel      string  `json:"white_user_level"`
	WhiteUserType       int     `json:"white_user_type" binding:"required"` // 1 玩家 2 假人 3 AI
	WhiteUserEnter      bool    `json:"white_user_enter"`
	WhiteUserExtra      string  `json:"white_user_extra"`
	BusinessType        string  `json:"business_type" binding:"required"` //业务类型
}

type ruleRequest struct {
	BoardSize     int `json:"board_size" binding:"required"` //9 13 19
	WinCaptured   int `json:"win_captured"`
	Type          int `json:"type" binding:"required"` // 1 吃子 2 围地
	HandicapCount int `json:"handicap_count"`          //让子 0 不让 1 让先 其他让子数
}

type studentHistoryRequest struct {
	StudentID uint  `json:"student_id"`
	StartDate int64 `json:"start_date"`
	EndDate   int64 `json:"end_date"`
}
