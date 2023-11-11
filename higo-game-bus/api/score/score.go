package score

import (
	"errors"
	"github.com/imroc/req/v3"
	"higo-game-bus/config"
	"higo-game-bus/request"
	"higo-game-bus/responses"
)

func AnalysisScoreApi(analysisRequest *request.AnalysisScoreRequest) (*responses.AnalysisScoreResult, error) {
	var response responses.AnalysisScoreResult
	resp, err := req.
		DevMode().
		R().
		SetBody(analysisRequest).
		SetSuccessResult(&response).
		Post("http://" + config.Conf.ThirdService.GoScore + "/api/katago-analysis/analysis")
	if err != nil {
		return nil, err
	}
	if resp.IsSuccessState() {
		return &response, nil
	}
	return nil, errors.New("请求错误")
}
