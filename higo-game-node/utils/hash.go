package utils

import (
	"crypto/md5"
	"encoding/hex"
	"higo-game-node/wq"
	"strconv"
)

func PointerHash(pointer [][]wq.Colour) string {
	str := ""
	for x := range pointer {
		for y := range pointer[x] {
			str += strconv.Itoa(int(pointer[x][y]))
		}
	}
	return getStringMD5(str)
}

func getStringMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
