package node

import (
	"game-lb/config"
	"game-lb/request"
	"testing"
)

func TestInfo(t *testing.T) {
	config.InitConf()
	Init()
	_, e := Info(&request.InfoRequest{
		BaseRequest: request.BaseRequest{
			GameId:   130090,
			UserId:   "3892",
			UserHash: "3892:1",
		},
	})
	if e != nil {
		t.Error(e)
	}
}
func TestBEnter(t *testing.T) {
	config.InitConf()
	Init()
	e := Enter(&request.EnterRequest{
		BaseRequest: request.BaseRequest{
			GameId:   130090,
			UserId:   "3892",
			UserHash: "3892:1",
		},
	})
	if e != nil {
		t.Error(e)
	}
}

func TestWEnter(t *testing.T) {
	config.InitConf()
	Init()
	e := Enter(&request.EnterRequest{
		BaseRequest: request.BaseRequest{
			GameId:   130090,
			UserId:   "66689",
			UserHash: "66689:1",
		},
	})
	if e != nil {
		t.Error(e)
	}
}

func TestBMove(t *testing.T) {
	config.InitConf()
	Init()
	_, e := Move(&request.MoveRequest{
		BaseRequest: request.BaseRequest{
			GameId:   130090,
			UserId:   "3892",
			UserHash: "3892:1",
		},
		X: 3,
		Y: 2,
		C: 1,
	})
	if e != nil {
		t.Error(e)
	}
}

func TestWMove(t *testing.T) {
	config.InitConf()
	Init()
	_, e := Move(&request.MoveRequest{
		BaseRequest: request.BaseRequest{
			GameId:   130090,
			UserId:   "66689",
			UserHash: "66689:1",
		},
		X: 4,
		Y: 2,
		C: 2,
	})
	if e != nil {
		t.Error(e)
	}
}

func TestInit(t *testing.T) {
	config.InitConf()
	Init()
	e := InitInfo(&request.InitRequest{
		GameId: 130090,
	})
	if e != nil {
		t.Error(e)
	}
}
