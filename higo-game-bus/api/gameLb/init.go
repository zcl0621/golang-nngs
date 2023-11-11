package gameLb

import (
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"higo-game-bus/config"
	"higo-game-bus/responses"
)

type InitRequest struct {
	GameId uint `json:"game_id"`
}

func InitInfoApi(request *InitRequest) error {
	var res responses.StandardResponse
	resp, err := req.
		R().
		SetBody(&request).
		SetSuccessResult(&res).
		Post(fmt.Sprintf("http://%s%s", config.Conf.ThirdService.GameLbService, "/api/v3/game-service/inner/init"))
	if err != nil {
		return err
	}
	if resp.IsSuccessState() {
		if res.Code == 0 {
			return nil
		} else {
			return errors.New(res.Msg)
		}
	}
	return errors.New("请求错误")
}
