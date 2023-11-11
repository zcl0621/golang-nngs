package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"higo-game-node/api"
	"higo-game-node/config"
	"higo-game-node/logger"
	"higo-game-node/responses"
	"time"
)

type GroupPublishData struct {
	GroupId string          `json:"group_id"` //消息组
	Message string          `json:"message"`  //消息数据 json字符串
	Index   int             `json:"index"`    //消息索引 提交消息时自增
	Ctx     context.Context `json:"-"`
}

type UserPublishData struct {
	UserId  string          `json:"user_id"` //用户id
	Message string          `json:"message"` //消息数据 json字符串
	Index   int             `json:"index"`   //消息索引 提交消息时自增
	Ctx     context.Context `json:"-"`
}

func groupPublishApi(request *GroupPublishData) error {
	r, _ := json.Marshal(request)
	var response responses.StandardResponse
	resp, err := api.ReqClint.
		R().
		SetBody(r).
		SetSuccessResult(&response).
		Post(fmt.Sprintf("http://%s%s", config.Conf.ThirdService.WsService, "/api/inner/ws-gateway-service/group-publish"))
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

func userPublishApi(request *UserPublishData) error {
	r, _ := json.Marshal(request)
	var response responses.StandardResponse
	resp, err := api.ReqClint.
		R().
		SetBody(r).
		SetSuccessResult(&response).
		Post(fmt.Sprintf("http://%s%s", config.Conf.ThirdService.WsService, "/api/inner/ws-gateway-service/user-publish"))
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

var UserPublishApiChan chan *UserPublishData
var GroupPublishApiChan chan *GroupPublishData

func StartWsPushChan() {
	UserPublishApiChan = make(chan *UserPublishData, 1024)
	GroupPublishApiChan = make(chan *GroupPublishData, 1024)
	go startUserPublish()
	go startGroupPublish()
}

func startUserPublish() {
	for {
		r := <-UserPublishApiChan
		if r.Ctx == nil {
			if e := userPublishApi(r); e != nil {
				logger.Logger("StartWsPushChan.UserPublishApi", logger.ERROR, e, fmt.Sprintf("WS消息推送失败 %v e %s", r, e.Error()))
			}
			continue
		}
		go func(d *UserPublishData) {
			ticker := time.NewTicker(time.Second * 5)
			defer ticker.Stop()
			select {
			case <-d.Ctx.Done():
				UserPublishApiChan <- &UserPublishData{
					UserId:  d.UserId,
					Message: d.Message,
					Index:   d.Index,
				}
				return
			case <-ticker.C:
				logger.Logger("StartWsPushChan.UserPublishApi", logger.ERROR, nil, fmt.Sprintf("WS消息推送超时 %v", d))
				return
			}
		}(r)
	}
}

func startGroupPublish() {
	for {
		r := <-GroupPublishApiChan
		if r.Ctx == nil {
			if e := groupPublishApi(r); e != nil {
				logger.Logger("StartWsPushChan.GroupPublishApi", logger.ERROR, e, fmt.Sprintf("WS消息推送失败 %v e %s", r, e.Error()))
			}
			continue
		}
		go func(d *GroupPublishData) {
			ticker := time.NewTicker(time.Second * 5)
			defer ticker.Stop()
			select {
			case <-d.Ctx.Done():
				GroupPublishApiChan <- &GroupPublishData{
					GroupId: d.GroupId,
					Message: d.Message,
					Index:   d.Index,
					Ctx:     nil,
				}
				return
			case <-ticker.C:
				logger.Logger("StartWsPushChan.GroupPublishApi", logger.ERROR, nil, fmt.Sprintf("WS消息推送超时 %v", d))
				return
			}
		}(r)
	}
}
