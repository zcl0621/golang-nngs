package public

import "higo-game-bus/responses"

type myBattleRequest struct {
	UserId         string `json:"user_id" form:"user_id"`
	BusinessType   string `json:"business_type" form:"business_type"`
	StartAtBegin   int64  `json:"start_at_begin" form:"start_at_begin"`     //筛选开始时间 区间A
	StartAtEnd     int64  `json:"start_at_end" form:"start_at_end"`         //筛选开始时间 区间B
	UserName       string `json:"user_name" form:"user_name"`               //用户名 手机号
	UserAccount    string `json:"user_account" form:"user_account"`         //学习卡号
	UserNickName   string `json:"user_nick_name" form:"user_nick_name"`     //用户昵称
	UserActualName string `json:"user_actual_name" form:"user_actual_name"` //用户真实姓名
	responses.PageRequest
}

type gameRequest struct {
	Type         string `json:"type" form:"type"` //territory captured
	UserId       string `json:"user_id" form:"user_id"`
	BusinessType string `json:"business_type" form:"business_type"`
	StartAtBegin int64  `json:"start_at_begin" form:"start_at_begin"` //筛选开始时间 区间A
	StartAtEnd   int64  `json:"start_at_end" form:"start_at_end"`     //筛选开始时间 区间B
	IsShow       bool   `json:"-"`
	NeedCount    bool   `json:"-"`
	responses.PageRequest
}

type GameInfoRequest struct {
	GameId uint `json:"game_id" form:"game_id" binding:"required"`
}

type businessViewRequest struct {
	StartTime         int64  `json:"start_time" binding:"required"`
	EndTime           int64  `json:"end_time" binding:"required"`
	TemporaryTestData string `json:"temporary_test_data"`
}
