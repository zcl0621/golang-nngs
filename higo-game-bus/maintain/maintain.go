package maintain

import "higo-game-bus/redisUtils"

var MaintaniGameKey = "higo:game:maintain"
var IsMaintain = false

func MaintainGame() error {
	IsMaintain = true
	return redisUtils.Set(MaintaniGameKey, []byte("1"), 24*3600)
}

func UnMaintainGame() error {
	IsMaintain = false
	has, e := redisUtils.Exists(MaintaniGameKey)
	if e != nil {
		return e
	}
	if has {
		return redisUtils.Del(MaintaniGameKey)
	}
	return nil
}

func CheckMaintain() bool {
	IsMaintain = false
	has, e := redisUtils.Exists(MaintaniGameKey)
	if e != nil {
		return false
	}
	if has {
		IsMaintain = true
		return true
	}
	return false
}
