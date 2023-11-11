package ai

import (
	"errors"
	"higo-game-node/api"
	"higo-game-node/config"
	"strconv"
	"strings"
)

type AiMoveResult struct {
	AiMove string `json:"ai_move"`
}

type AiMoveRequest struct {
	Level     uint    `json:"level"`
	StepTime  float32 `json:"stepTime"`
	BoardSize uint    `json:"boardSize"`
	Type      string  `json:"type"`
	SGF       string  `json:"sgf"`
}

func MoveApi(request *AiMoveRequest) (string, error) {
	responses := &AiMoveResult{}
	resp, err := api.ReqClint.R().
		SetBody(request).
		Post("http://" + config.Conf.ThirdService.AiAgent + "/ai/genmove")
	if err != nil {
		return "", err
	}
	if resp.IsSuccessState() {
		e := resp.Into(&responses)
		if e != nil {
			return "", e
		}
		return responses.AiMove, nil
	}
	return "", errors.New("ai error")
}

func AiMoveToXY(aiMove string, boardSize int) (int, int) {
	var aiMoveXY map[string]int
	aiMoveXY = make(map[string]int)
	aiMoveXY["X"] = strings.Index("ABCDEFGHJKLMNOPQRST", string(aiMove[0]))
	y, _ := strconv.Atoi(string(aiMove[1:]))
	aiMoveXY["Y"] = int(boardSize) - y
	return aiMoveXY["X"], aiMoveXY["Y"]
}
