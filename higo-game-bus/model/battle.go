package model

import (
	"errors"
	"gorm.io/gorm"
	"higo-game-bus/database"
	"higo-game-bus/exception"
)

type Battle struct {
	gorm.Model
	GameId                 uint   `json:"game_id" gorm:"index"`
	BoardSize              int    `json:"board_size"`
	Type                   string `json:"type"` // territory  captured
	UserId                 string `json:"user_id" gorm:"index"`
	UserName               string `json:"user_name"`
	UserAccount            string `json:"user_account"`
	UserNickName           string `json:"user_nick_name"`
	UserActualName         string `json:"user_actual_name"`
	UserAvatar             string `json:"user_avatar"`
	UserLevel              string `json:"user_level"`
	UserSide               int    `json:"user_side"`  // 1 黑 2 白
	UserWin                int    `json:"user_win"`   // 1 胜 2 负 3 和 4 弃权
	WinResult              string `json:"win_result"` // 输赢原因 B+2.5 黑胜2.5目 W+2.5 白胜2.5目 B+R 黑胜白认输 B+T 黑胜白超时 W+R 白胜黑认输 W+T 白胜黑超时 B+C1 黑胜吃1子 W+C1 白胜吃一子 B+L 黑胜白未参加 W+L  白胜黑未参加 B+A 黑胜仲裁 W+A 白胜仲裁 Draw 和 Abstain 弃权
	OpponentUserId         string `json:"opponent_user_id"`
	OpponentUserName       string `json:"opponent_user_name"`
	OpponentUserAccount    string `json:"opponent_user_account"`
	OpponentUserNickName   string `json:"opponent_user_nick_name"`
	OpponentUserActualName string `json:"opponent_user_actual_name"`
	OpponentUserAvatar     string `json:"opponent_user_avatar"`
	OpponentUserLevel      string `json:"opponent_user_level"`
	OpponentUserSide       int    `json:"opponent_user_side"`         // 1 黑 2 白
	BusinessType           string `json:"business_type" gorm:"index"` //业务类型
	StartedAt              int64  `json:"started_at"`
	EndedAt                int64  `json:"ended_at"`
	IsShow                 bool   `json:"is_show" gorm:"default:false;index"`
}

// SelectOneBattle
// swagger:ignore @title    查询数据
// swagger:ignore @param     tx        *gorm.DB         "事务db连接 可为nil"
// swagger:ignore @param     query        *Battle         "查询结构体 不可为nil"
// swagger:ignore @param 	 otherQuery   []func "其他查询条件"
// swagger:ignore @return    data        *Battle        "查询的结果"
func SelectOneBattle(tx *gorm.DB, query *Battle, otherQuery ...func(db *gorm.DB) *gorm.DB) *Battle {
	if query == nil && otherQuery == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("查询参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("SelectOneBattle 参数错误")))
	}
	if tx == nil {
		tx = database.GetInstance()
	}
	tx = tx.Model(&Battle{})
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
	var data Battle
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
		SetOriginalError(errors.New("SelectOneBattle 查询异常")))
}

// SelectBattles
// swagger:ignore @title    查询数据
// swagger:ignore @param     tx        *gorm.DB         "事务db连接 可为nil"
// swagger:ignore @param     query        *Battle         "查询结构体 不可为nil"
// swagger:ignore @param 	 otherQuery   []func "其他查询条件"
// swagger:ignore @return    data        []Battle        "查询的结果"
func SelectBattles(tx *gorm.DB, query *Battle, otherQuery ...func(db *gorm.DB) *gorm.DB) []Battle {
	if query == nil && otherQuery == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("查询参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("SelectBattle 参数错误")))
	}
	if tx == nil {
		tx = database.GetInstance()
	}
	tx = tx.Model(&Battle{})
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
	var data []Battle
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
		SetOriginalError(errors.New("SelectBattle 查询异常")))
}

// ScanBattles
// swagger:ignore @title    查询数据
// swagger:ignore @param     tx        *gorm.DB         "事务db连接 可为nil"
// swagger:ignore @param     query        *Battle         "查询结构体 不可为nil"
// swagger:ignore @param 	 otherQuery   func "其他查询条件"
func ScanBattles(tx *gorm.DB, query *Battle, otherQuery func(db *gorm.DB) *gorm.DB) {
	if query == nil && otherQuery == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("查询参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("ScanBattle 参数错误")))
	}
	if tx == nil {
		tx = database.GetInstance()
	}
	tx = tx.Model(&Battle{})
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
		SetOriginalError(errors.New("ScanBattle 查询异常")))
}

// CreateBattles
// swagger:ignore @title    创建数据
// swagger:ignore @param     tx        *gorm.DB         "事务db连接 可为nil"
// swagger:ignore @param     create        ...*Battle         "创建结构体 不可为nil"
func CreateBattles(tx *gorm.DB, create ...*Battle) {
	if create == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("创建参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("CreateBattles 参数错误")))
	}
	if tx == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("CreateBattles 参数错误")))
	}
	tx = tx.Model(&Battle{})
	var querySet = tx.Model(&Battle{}).Create(create)
	if querySet.Error == nil {
		return
	}
	panic(exception.StandardRuntimeBadError().
		SetOutPutMessage("创建异常").
		SetErrorCode(2).
		SetParameter(create).
		SetOriginalError(errors.New("CreateBattles 创建异常")))
}

// UpdateBattles
// swagger:ignore @title    更新数据
// swagger:ignore @param     tx        *gorm.DB         "事务db连接 可为nil"
// swagger:ignore @param     query        *Battle         "查询结构体 不可为nil"
// swagger:ignore @param     update        *Battle         "更新结构体 不可为nil"
// swagger:ignore @param 	 otherQuery   []func "其他查询条件"
func UpdateBattles(tx *gorm.DB, query *Battle, update *Battle, otherQuery ...func(db *gorm.DB) *gorm.DB) {
	if query == nil && otherQuery == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("更新参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("UpdateBattles 参数错误")))
	}
	if update == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("UpdateBattles 参数错误")))
	}
	if tx == nil {
		tx = database.GetInstance()
	}
	tx = tx.Model(&Battle{})
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
		SetOriginalError(errors.New("UpdateBattles 更新异常")))
}

// DeleteBattles
// swagger:ignore @title    删除数据
// swagger:ignore @param     tx        *gorm.DB         "事务db连接 可为nil"
// swagger:ignore @param     query        *Battle         "查询结构体 不可为nil"
// swagger:ignore @param 	 otherQuery   func "其他查询条件"
func DeleteBattles(tx *gorm.DB, query *Battle, otherQuery ...func(db *gorm.DB) *gorm.DB) {
	if query == nil && otherQuery == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("删除参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("DeleteBattles 参数错误")))
	}
	if tx == nil {
		panic(exception.StandardRuntimeBadError().
			SetOutPutMessage("参数错误").
			SetErrorCode(2).
			SetOriginalError(errors.New("DeleteBattles 参数错误")))
	}
	tx = tx.Model(&Battle{})
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
	var querySet = tx.Delete(&Battle{})
	if querySet.Error == nil {
		return
	}
	panic(exception.StandardRuntimeBadError().
		SetOutPutMessage("删除异常").
		SetErrorCode(2).
		SetParameter(query).
		SetOriginalError(errors.New("DeleteBattles 删除异常")))
}
