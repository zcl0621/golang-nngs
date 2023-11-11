package model

import (
	"errors"
	"gorm.io/gorm"
	"higo-game-bus/database"
	"higo-game-bus/exception"
)

// GameBusinessTypeMap 对弈类型
var GameBusinessTypeMap = map[string]string{
	"tournament":    "聂道赛事",
	"ai-tournament": "AI对弈",
	"ai-stage":      "课后对弈",
	"room_solo":     "好友约战",
	"auto_pair":     "排位赛",
	"rank_level":    "定级赛",
	"smart_board":   "智能棋盘",
	"season-rank":   "排位赛",
	"evaluation":    "评测",
	"unit":          "单元测",
}

var HUMANUSER = 1
var FAKEUSER = 2
var AIUSER = 3

type ExtraData struct {
	AiAgentLevel    uint    `json:"ai_agent_level"`
	AiAgentStepTime float32 `json:"ai_agent_step_time"`
}

type Game struct {
	gorm.Model
	BoardSize            int     `json:"board_size" gorm:"default:19"` //9 13 19
	SGF                  string  `json:"sgf"`
	CanStartTime         int64   `json:"can_start_time" gorm:"default:0"`
	NotEnterTime         int64   `json:"not_enter_time" gorm:"default:0"`   // 一方未进入或两方都未进入 强制结束时间
	EnableMoveTime       int     `json:"enable_move_time" gorm:"default:0"` // 是否开启第一步倒计时 0 否 1 是
	IsStart              bool    `json:"is_start"`
	StartTime            int64   `json:"start_time" gorm:"default:0"`
	IsEnd                bool    `json:"is_end"`
	EndTime              int64   `json:"end_time" gorm:"default:0"`
	WinCaptured          int     `json:"win_captured" gorm:"default:0"`          // 赢棋子数
	MaxStep              int     `json:"max_step" gorm:"default:0"`              // 最大步数
	KM                   float64 `json:"km" gorm:"default:0;type:decimal(10,2)"` //贴目
	TerritoryStep        int     `json:"territory_step" gorm:"default:0"`
	BlackReturnStone     float64 `json:"black_return_stone" gorm:"default:0;type:decimal(10,2)"` //黑返子数
	Win                  int     `json:"win" gorm:"default:0"`                                   // 0 无 1 黑 2 白 3 和 4 弃权
	WinResult            string  `json:"win_result"`                                             // 输赢原因 B+2.5 黑胜2.5目 W+2.5 白胜2.5目 B+R 黑胜白认输 B+T 黑胜白超时 W+R 白胜黑认输 W+T 白胜黑超时 B+C1 黑胜吃1子 W+C1 白胜吃一子 B+L 黑胜白未参加 W+L  白胜黑未参加 B+A 黑胜仲裁 W+A 白胜仲裁 Draw 和 Abstain 弃权
	Step                 int     `json:"step" gorm:"default:0"`
	BlackTime            int     `json:"black_time" gorm:"default:0"`
	WhiteTime            int     `json:"white_time" gorm:"default:0"`
	BlackByoYomi         int     `json:"black_byo_yomi" gorm:"default:0"`
	WhiteByoYomi         int     `json:"white_byo_yomi" gorm:"default:0"`
	BlackByoYomiTime     int     `json:"black_byo_yomi_time" gorm:"default:0"`
	WhiteByoYomiTime     int     `json:"white_byo_yomi_time" gorm:"default:0"`
	LeftBlackTime        int     `json:"left_black_time" gorm:"default:0"`
	LeftWhiteTime        int     `json:"left_white_time" gorm:"default:0"`
	LeftBlackByoYomiTime int     `json:"left_black_byo_yomi_time" gorm:"default:0"`
	LeftWhiteByoYomiTime int     `json:"left_white_byo_yomi_time" gorm:"default:0"`
	LeftWhiteByoYomi     int     `json:"left_white_byo_yomi" gorm:"default:0"`
	LeftBlackByoYomi     int     `json:"left_black_byo_yomi" gorm:"default:0"`
	BlackUserId          string  `json:"black_user_id"`
	BlackUserName        string  `json:"black_user_name"`
	BlackUserAccount     string  `json:"black_user_account"`
	BlackUserNickName    string  `json:"black_user_nick_name"`
	BlackUserActualName  string  `json:"black_user_actual_name"`
	BlackUserAvatar      string  `json:"black_user_avatar"`
	BlackUserLevel       string  `json:"black_user_level"`
	BlackUserType        int     `json:"black_user_type"`  // 1 玩家 2 假人 3 AI
	BlackUserExtra       string  `json:"black_user_extra"` // 附加信息
	BlackUserEnter       bool    `json:"black_enter"`
	BlackUserHash        string  `json:"black_user_hash" gorm:"index"`
	WhiteUserId          string  `json:"white_user_id"`
	WhiteUserName        string  `json:"white_user_name"`
	WhiteUserAccount     string  `json:"white_user_account"`
	WhiteUserNickName    string  `json:"white_user_nick_name"`
	WhiteUserActualName  string  `json:"white_user_actual_name"`
	WhiteUserAvatar      string  `json:"white_user_avatar"`
	WhiteUserLevel       string  `json:"white_user_level"`
	WhiteUserType        int     `json:"white_user_type"` // 1 玩家 2 假人 3 AI
	WhiteUserHash        string  `json:"white_user_hash" gorm:"index"`
	WhiteUserEnter       bool    `json:"white_user_enter"`
	WhiteUserExtra       string  `json:"white_user_extra"`
	BusinessType         string  `json:"business_type" gorm:"index"` //业务类型
	ViewCount            int     `json:"view_count" gorm:"default:0"`

	BlackCaptured int     `json:"black_captured" gorm:"default:0"`
	WhiteCaptured int     `json:"white_captured" gorm:"default:0"`
	BlackScore    float64 `json:"black_score" gorm:"default:0;type:decimal(10,2)"`
	WhiteScore    float64 `json:"white_score" gorm:"default:0;type:decimal(10,2)"`
	IsShow        bool    `json:"is_show" gorm:"default:false;index"`
}

// SelectOneGame
// swagger:ignore @title    查询数据
// swagger:ignore @param     tx        *gorm.DB         "事务db连接 可为nil"
// swagger:ignore @param     query        *Game         "查询结构体 不可为nil"
// swagger:ignore @param 	 otherQuery   []func "其他查询条件"
// swagger:ignore @return    data        *Game        "查询的结果"
func SelectOneGame(tx *gorm.DB, query *Game, otherQuery ...func(db *gorm.DB) *gorm.DB) *Game {
	if query == nil && otherQuery == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("查询参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("SelectOneGame 参数错误")))
	}
	if tx == nil {
		tx = database.GetInstance()
	}
	tx = tx.Model(&Game{})
	if query != nil {
		tx = tx.Where(query)
	}
	if otherQuery != nil {
		for i := range otherQuery {
			if otherQuery[i] != nil {
				tx = otherQuery[i](tx)
			}
		}
	}
	var data Game
	var querySet = tx.First(&data)
	if querySet.Error == nil {
		return &data
	}
	if querySet.RowsAffected == 0 {
		return nil
	}
	panic(exception.StandardRuntimeBadError().
		SetOutPutMessage("查询异常").
		SetErrorCode(2).
		SetParameter(query).
		SetOriginalError(errors.New("SelectOneGame 查询异常")))
}

// SelectGames
// swagger:ignore @title    查询数据
// swagger:ignore @param     tx        *gorm.DB         "事务db连接 可为nil"
// swagger:ignore @param     query        *Game         "查询结构体 不可为nil"
// swagger:ignore @param 	 otherQuery   []func "其他查询条件"
// swagger:ignore @return    data        []Game        "查询的结果"
func SelectGames(tx *gorm.DB, query *Game, otherQuery ...func(db *gorm.DB) *gorm.DB) []Game {
	if query == nil && otherQuery == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("查询参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("SelectGame 参数错误")))
	}
	if tx == nil {
		tx = database.GetInstance()
	}
	tx = tx.Model(&Game{})
	if query != nil {
		tx = tx.Where(query)
	}
	if otherQuery != nil {
		for i := range otherQuery {
			if otherQuery[i] != nil {
				tx = otherQuery[i](tx)
			}
		}
	}
	var data []Game
	var querySet = tx.Find(&data)
	if querySet.Error == nil {
		return data
	}
	if querySet.RowsAffected == 0 {
		return nil
	}
	panic(exception.StandardRuntimeBadError().
		SetOutPutMessage("查询异常").
		SetErrorCode(2).
		SetParameter(query).
		SetOriginalError(errors.New("SelectGame 查询异常")))
}

// ScanGames
// swagger:ignore @title    查询数据
// swagger:ignore @param     tx        *gorm.DB         "事务db连接 可为nil"
// swagger:ignore @param     query        *Game         "查询结构体 不可为nil"
// swagger:ignore @param 	 otherQuery   func "其他查询条件"
func ScanGames(tx *gorm.DB, query *Game, otherQuery func(db *gorm.DB) *gorm.DB) {
	if query == nil && otherQuery == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("查询参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("ScanGame 参数错误")))
	}
	if tx == nil {
		tx = database.GetInstance()
	}
	tx = tx.Model(&Game{})
	if query != nil {
		tx = tx.Where(query)
	}
	querySet := otherQuery(tx)
	if querySet.Error == nil {
		return
	}
	panic(exception.StandardRuntimeBadError().
		SetOutPutMessage("查询异常").
		SetErrorCode(2).
		SetParameter(query).
		SetOriginalError(errors.New("ScanGame 查询异常")))
}

// CreateGames
// swagger:ignore @title    创建数据
// swagger:ignore @param     tx        *gorm.DB         "事务db连接 可为nil"
// swagger:ignore @param     create        ...*Game         "创建结构体 不可为nil"
func CreateGames(tx *gorm.DB, create ...*Game) {
	if create == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("创建参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("CreateGames 参数错误")))
	}
	if tx == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("CreateGames 参数错误")))
	}
	tx = tx.Model(&Game{})
	var querySet = tx.Model(&Game{}).Create(create)
	if querySet.Error == nil {
		return
	}
	panic(exception.StandardRuntimeBadError().
		SetOutPutMessage("创建异常").
		SetErrorCode(2).
		SetParameter(create).
		SetOriginalError(errors.New("CreateGames 创建异常")))
}

// UpdateGames
// swagger:ignore @title    更新数据
// swagger:ignore @param     tx        *gorm.DB         "事务db连接 可为nil"
// swagger:ignore @param     query        *Game         "查询结构体 不可为nil"
// swagger:ignore @param     update        *Game         "更新结构体 不可为nil"
// swagger:ignore @param 	 otherQuery   []func "其他查询条件"
func UpdateGames(tx *gorm.DB, query *Game, update *Game, otherQuery ...func(db *gorm.DB) *gorm.DB) {
	if query == nil && otherQuery == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("更新参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("UpdateGames 参数错误")))
	}
	if update == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("UpdateGames 参数错误")))
	}
	if tx == nil {
		tx = database.GetInstance()
	}
	tx = tx.Model(&Game{})
	if query != nil {
		tx = tx.Where(query)
	}
	if otherQuery != nil {
		for i := range otherQuery {
			if otherQuery[i] != nil {
				tx = otherQuery[i](tx)
			}
		}
	}
	var querySet = tx.Updates(update)
	if querySet.Error == nil {
		return
	}
	panic(exception.StandardRuntimeBadError().
		SetOutPutMessage("更新异常").
		SetErrorCode(2).
		SetParameter(query).
		SetOriginalError(errors.New("UpdateGames 更新异常")))
}

// DeleteGames
// swagger:ignore @title    删除数据
// swagger:ignore @param     tx        *gorm.DB         "事务db连接 可为nil"
// swagger:ignore @param     query        *Game         "查询结构体 不可为nil"
// swagger:ignore @param 	 otherQuery   func "其他查询条件"
func DeleteGames(tx *gorm.DB, query *Game, otherQuery ...func(db *gorm.DB) *gorm.DB) {
	if query == nil && otherQuery == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("删除参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("DeleteGames 参数错误")))
	}
	if tx == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("DeleteGames 参数错误")))
	}
	tx = tx.Model(&Game{})
	if query != nil {
		tx = tx.Where(query)
	}
	if otherQuery != nil {
		for i := range otherQuery {
			if otherQuery[i] != nil {
				tx = otherQuery[i](tx)
			}
		}
	}
	var querySet = tx.Delete(&Game{})
	if querySet.Error == nil {
		return
	}
	panic(exception.StandardRuntimeBadError().
		SetOutPutMessage("删除异常").
		SetErrorCode(2).
		SetParameter(query).
		SetOriginalError(errors.New("DeleteGames 删除异常")))
}
