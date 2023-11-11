package gameBus

import (
	"errors"
	"higo-game-node/api"
	"higo-game-node/config"
	"higo-game-node/responses"
	"higo-game-node/utils"
	"higo-game-node/wq"
)

type AnalysisScoreRequest struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

type AnalysisScoreResultRequest struct {
	Id string `json:"id"`
}

type AnalysisScoreData struct {
	Id        string           `json:"id"`
	Ownership []float64        `json:"ownership"`
	RootInfo  AnalysisRootInfo `json:"rootInfo"`
}

type AnalysisRootInfo struct {
	CurrentPlayer string  `json:"currentPlayer"`
	ScoreLead     float64 `json:"scoreLead"`
	ScoreSelfplay float64 `json:"scoreSelfplay"`
	ScoreStdev    float64 `json:"scoreStdev"`
	SymHash       string  `json:"symHash"`
	ThisHash      string  `json:"thisHash"`
	Utility       float64 `json:"utility"`
	Visits        int     `json:"visits"`
	Winrate       float64 `json:"winrate"`
}

var SCOREREDISRESULTKEY = "higo:game:bus:score:result:%s"

type AnalysisScoreResult struct {
	Id      string `json:"id"`
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

func StartAnalysisScore(data *wq.KataGoAnalysis) error {
	var request AnalysisScoreRequest
	request.Data = utils.ZipString(data.ToString())
	request.Id = data.Id
	var response responses.StandardResponse
	resp, err := api.ReqClint.
		R().
		SetBody(request).
		SetSuccessResult(&response).
		Post("http://" + config.Conf.ThirdService.GameBusService + "/api/v3/game-service/inner/score/start")
	if err != nil {
		return err
	}
	if resp.IsSuccessState() {
		if response.Code == 0 {
			return nil
		} else {
			return errors.New(response.Msg)
		}
	}
	return errors.New("请求错误")
}

func GetAnalysisScore(ownership *[]float64, blackReturnStones float64) (bScore float64, wScore float64, endScore float64, controversyCount int) {
	for _, value := range *ownership {
		if value >= 0.5 {
			bScore++
		} else if value < 0.5 && value >= 0.25 {
			bScore += 0.5
		} else if value > -0.5 && value <= -0.25 {
			wScore += 0.5
		} else if value <= -0.5 {
			wScore++
		} else {
			controversyCount++
		}
	}
	bScore += float64(controversyCount) / 2
	wScore += float64(controversyCount) / 2
	//controversyCount = controversyCount % 2
	//if currentPlayer == "B" {
	//	bScore += float64(controversyCount)
	//} else {
	//	wScore += float64(controversyCount)
	//}
	endScore = bScore - blackReturnStones - wScore
	return
}
