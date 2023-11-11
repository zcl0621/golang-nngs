package notifyStruct

import (
	"encoding/json"
	"fmt"
	"higo-game-node/utils"
	"higo-game-node/wq"
)

type GameMoveMsg struct {
	GameId      uint      `json:"game_id"`
	X           int       `json:"x"`
	Y           int       `json:"y"`
	C           wq.Colour `json:"c"`
	Turn        wq.Colour `json:"turn"`
	Step        int       `json:"step"`
	PointerHash string    `json:"pointer_hash"` // 盘面hash
}

type WSGameMoveMsg struct {
	MessageType string      `json:"message_type"`
	Data        GameMoveMsg `json:"data"`
}

func MakeWsGameMoveMsg(data *GameMoveMsg) string {
	msg := WSGameMoveMsg{MessageType: "move", Data: *data}
	d, _ := json.Marshal(msg)
	return utils.ZipString(fmt.Sprintf("%s", d))
}

type GameEndMsg struct {
	GameId           uint           `json:"game_id"`
	Win              int            `json:"win"`
	WinResult        string         `json:"win_result"`
	From             string         `json:"from"`
	Step             int            `json:"step"`
	BCaptured        int            `json:"b_captured"`
	WCaptured        int            `json:"w_captured"`
	BScore           float64        `json:"b_score"`
	WScore           float64        `json:"w_score"`
	ControversyCount int            `json:"controversy_count"`
	OwnerShip        []wq.OwnerShip `json:"owner_ship"`
}

type WSGameEndMsg struct {
	MessageType string     `json:"message_type"`
	Data        GameEndMsg `json:"data"`
}

func MakeWsGameEndMsg(data *GameEndMsg) string {
	msg := WSGameEndMsg{MessageType: "end", Data: *data}
	d, _ := json.Marshal(msg)
	return utils.ZipString(fmt.Sprintf("%s", d))
}

type GamePassMsg struct {
	GameId      uint      `json:"game_id"`
	C           wq.Colour `json:"c"`
	Step        int       `json:"step"`
	PointerHash string    `json:"pointer_hash"`
}

type WSGamePassMsg struct {
	MessageType string      `json:"message_type"`
	Data        GamePassMsg `json:"data"`
}

func MakeWsGamePassMsg(data *GamePassMsg) string {
	msg := WSGamePassMsg{MessageType: "pass", Data: *data}
	d, _ := json.Marshal(msg)
	return utils.ZipString(fmt.Sprintf("%s", d))
}

type GameTimeMsg struct {
	GameId           uint      `json:"game_id"`
	C                wq.Colour `json:"c"`
	BlackMainTime    int       `json:"black_main_time"`
	BlackByoYomiTime int       `json:"black_byo_yomi_time"`
	BlackByoYomi     int       `json:"black_byo_yomi"`
	WhiteMainTime    int       `json:"white_main_time"`
	WhiteByoYomiTime int       `json:"white_byo_yomi_time"`
	WhiteByoYomi     int       `json:"white_byo_yomi"`
}

type WSGameTimeMsg struct {
	MessageType string      `json:"message_type"`
	Data        GameTimeMsg `json:"data"`
}

func MakeWsGameTimeMsg(data *GameTimeMsg) string {
	msg := WSGameTimeMsg{MessageType: "time", Data: *data}
	d, _ := json.Marshal(msg)
	return utils.ZipString(fmt.Sprintf("%s", d))
}

type AreaScoreTimeMsg struct {
	GameId        uint   `json:"game_id"`
	LeftTime      int    `json:"left_time"`
	ApplyUseId    string `json:"apply_use_id"`
	ApplyUserHash string `json:"apply_user_hash"`
	OpponentId    string `json:"opponent_id"`
	OpponentHash  string `json:"opponent_hash"`
	Status        int    `json:"status"` // 1 等待对手同意 2 对手拒绝 3 对手同意 4 对手超时
}

type WSAreaScoreTimeMsg struct {
	MessageType string           `json:"message_type"`
	Data        AreaScoreTimeMsg `json:"data"`
}

func MakeWsAreaScoreTimeMsg(data *AreaScoreTimeMsg) string {
	msg := WSAreaScoreTimeMsg{MessageType: "area_score_time", Data: *data}
	d, _ := json.Marshal(msg)
	return utils.ZipString(fmt.Sprintf("%s", d))
}

type UserStatusMsg struct {
	GameId       uint   `json:"game_id"`
	UserId       string `json:"user_id"`
	UserHash     string `json:"user_hash"`
	OnlineStatus bool   `json:"online_status"`
	EnterStatus  bool   `json:"enter_status"`
}

type WSUserStatusMsg struct {
	MessageType string        `json:"message_type"`
	Data        UserStatusMsg `json:"data"`
}

func MakeWsUserStatusMsg(data *UserStatusMsg) string {
	msg := WSUserStatusMsg{MessageType: "user_status", Data: *data}
	d, _ := json.Marshal(msg)
	return utils.ZipString(fmt.Sprintf("%s", d))
}

type SummationTimeMsg struct {
	GameId        uint   `json:"game_id"`
	LeftTime      int    `json:"left_time"`
	ApplyUseId    string `json:"apply_use_id"`
	ApplyUserHash string `json:"apply_user_hash"`
	OpponentId    string `json:"opponent_id"`
	OpponentHash  string `json:"opponent_hash"`
	Status        int    `json:"status"` // 1 等待对手同意 2 对手拒绝 3 对手同意 4 对手超时
}

type WSSummationTimeMsg struct {
	MessageType string           `json:"message_type"`
	Data        SummationTimeMsg `json:"data"`
}

func MakeWsSummationTimeMsg(data *SummationTimeMsg) string {
	msg := WSSummationTimeMsg{MessageType: "summation_time", Data: *data}
	d, _ := json.Marshal(msg)
	return utils.ZipString(fmt.Sprintf("%s", d))
}

type MoveTimeMsg struct {
	GameId   uint      `json:"game_id"`
	Color    wq.Colour `json:"color"`
	LeftTime int       `json:"left_time"`
}

type WSMoveTimeMsg struct {
	MessageType string      `json:"message_type"`
	Data        MoveTimeMsg `json:"data"`
}

func MakeWsMoveTimeMsg(data *MoveTimeMsg) string {
	msg := WSMoveTimeMsg{MessageType: "move_time", Data: *data}
	d, _ := json.Marshal(msg)
	return utils.ZipString(fmt.Sprintf("%s", d))
}

type ForceReloadSGFMsg struct {
	GameId uint   `json:"game_id"`
	SGF    string `json:"sgf"`
}

type WSForceReloadSGFMsg struct {
	MessageType string            `json:"message_type"`
	Data        ForceReloadSGFMsg `json:"data"`
}

func MakeWsForceReloadSGFMsg(data *ForceReloadSGFMsg) string {
	msg := WSForceReloadSGFMsg{MessageType: "force_reload", Data: *data}
	d, _ := json.Marshal(msg)
	return utils.ZipString(fmt.Sprintf("%s", d))
}

type StartByomiMsg struct {
	GameId uint      `json:"game_id"`
	C      wq.Colour `json:"c"`
}

type WSStartByomiMsg struct {
	MessageType string        `json:"message_type"`
	Data        StartByomiMsg `json:"data"`
}

func MakeWsStartByomiMsg(data *StartByomiMsg) string {
	msg := WSStartByomiMsg{MessageType: "start_byomi", Data: *data}
	d, _ := json.Marshal(msg)
	return utils.ZipString(fmt.Sprintf("%s", d))
}
