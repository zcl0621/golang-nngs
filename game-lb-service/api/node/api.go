package node

import (
	"context"
	"game-lb/request"
	"game-lb/responses"
	"time"
)

func InitInfo(req *request.InitRequest) (err error) {
	var resp string
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "InitInfo", *req, &resp)
	if err != nil {
		return err
	}
	return nil
}

// Info 获取对弈信息
func Info(req *request.InfoRequest) (*responses.InfoResponse, error) {
	var err error
	resp := &responses.InfoResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "Info", *req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SGF 获取对弈信息
func SGF(req *request.InfoRequest) (resp *responses.SGFResponse, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "SGFInfo", *req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Enter 进入对弈
func Enter(req *request.EnterRequest) (err error) {
	var resp string
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "Enter", *req, &resp)
	if err != nil {
		return err
	}
	return nil
}

// Move 落子
func Move(req *request.MoveRequest) (*responses.MoveResponse, error) {
	var err error
	resp := &responses.MoveResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "Move", *req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Pass 停一手
func Pass(req *request.PassRequest) (*responses.PassResponse, error) {
	var err error
	resp := &responses.PassResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "Pass", *req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Resign 认输
func Resign(req *request.ResignRequest) (*responses.EndResponse, error) {
	var err error
	resp := &responses.EndResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "Resign", *req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AreaScore 数目
func AreaScore(req *request.InfoRequest) (*responses.EndResponse, error) {
	var err error
	resp := &responses.EndResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "AreaScore", *req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ApplyAreaScore 申请数目
func ApplyAreaScore(req *request.InfoRequest) (err error) {
	var resp string
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "ApplyAreaScore", *req, &resp)
	if err != nil {
		return err
	}
	return nil
}

// AgreeAreaScore 同意数目
func AgreeAreaScore(req *request.InfoRequest) (*responses.EndResponse, error) {
	var err error
	resp := &responses.EndResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "AgreeAreaScore", *req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RejectAreaScore 拒绝数目
func RejectAreaScore(req *request.InfoRequest) (err error) {
	var resp string
	err = nodeXClient.Call(context.Background(), "RejectAreaScore", *req, &resp)
	if err != nil {
		return err
	}
	return nil
}

// OwnerShip 形势
func OwnerShip(req *request.OwnershipRequest) (*responses.OwnerShipResponse, error) {
	var err error
	resp := &responses.OwnerShipResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "Ownership", *req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CallEnd 结束
func CallEnd(req *request.CallEndRequest) (err error) {
	var resp string
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "CallEnd", *req, &resp)
	if err != nil {
		return err
	}
	return nil
}

// ApplySummation 申请和棋
func ApplySummation(req *request.InfoRequest) (err error) {
	var resp string
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "ApplySummation", *req, &resp)
	if err != nil {
		return err
	}
	return nil
}

// AgreeSummation 同意和棋
func AgreeSummation(req *request.InfoRequest) (*responses.EndResponse, error) {
	var err error
	resp := &responses.EndResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "AgreeSummation", *req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RejectSummation 拒绝和棋
func RejectSummation(req *request.InfoRequest) (err error) {
	var resp string
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "RejectSummation", *req, &resp)
	if err != nil {
		return err
	}
	return nil
}

// CanPlay 当前颜色
func CanPlay(req *request.InfoRequest) (*responses.CanPlayResponse, error) {
	var err error
	resp := &responses.CanPlayResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "CanPlay", *req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// InnerAreaScore 数目
func InnerAreaScore(req *request.OwnershipRequest) (*responses.InnerAreaScoreResponse, error) {
	var err error
	resp := &responses.InnerAreaScoreResponse{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "InnerAreaScore", *req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ForceReloadSGF 强制重新加载SGF
func ForceReloadSGF(req *request.ForceReloadSgfRequest) (err error) {
	var resp string
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = nodeXClient.Call(ctx, "ForceReloadSGF", *req, &resp)
	if err != nil {
		return
	}
	return
}
