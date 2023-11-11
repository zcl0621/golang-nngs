package play

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"higo-game-node/api/gameBus"
	"higo-game-node/cache"
	"higo-game-node/logger"
	"higo-game-node/redisUtils"
	"higo-game-node/utils"
	"higo-game-node/wq"
	"time"
)

type Arith struct{}

// SGFInfo 棋谱
func (t *Arith) SGFInfo(ctx context.Context, args InfoRequest, reply *SGFResponse) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("SGFInfo", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	e = sgfService(&args, reply)
	return e
}

// InitInfo 创建基础信息 内部使用
func (t *Arith) InitInfo(ctx context.Context, args InitRequest, reply *string) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("InitInfo", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	e = initService(&args)
	return e
}

// Info 基础信息
func (t *Arith) Info(ctx context.Context, args InfoRequest, reply *InfoResponse) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("Info", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	e = infoService(&args, reply)
	return e
}

// Enter 进入
func (t *Arith) Enter(ctx context.Context, args EnterRequest, reply *string) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("Enter", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	e = enterService(&args, d)
	return e
}

// Move 落子
func (t *Arith) Move(ctx context.Context, args MoveRequest, reply *MoveResponse) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("Move", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	if !d.Lock.TryLockWithTimeout(time.Second * 3) {
		e = errors.New("请求过于频繁")
		return e
	}
	e = moveService(&args, d, reply)
	d.Lock.Unlock()
	return e
}

// Pass 停一手
func (t *Arith) Pass(ctx context.Context, args PassRequest, reply *PassResponse) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("Pass", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	if !d.Lock.TryLockWithTimeout(time.Second * 3) {
		e = errors.New("请求过于频繁")
		return e
	}
	e = passService(&args, d, reply)
	d.Lock.Unlock()
	return e
}

// Resign 认输
func (t *Arith) Resign(ctx context.Context, args ResignRequest, reply *EndResponse) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("Resign", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	if !d.Lock.TryLockWithTimeout(time.Second * 3) {
		e = errors.New("请求过于频繁")
		return e
	}
	e = resignService(&args, d, reply)
	d.Lock.Unlock()
	return e
}

// AreaScore 数目
func (t *Arith) AreaScore(ctx context.Context, args InfoRequest, reply *EndResponse) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("AreaScore", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	if !d.Lock.TryLockWithTimeout(time.Second * 3) {
		e = errors.New("请求过于频繁")
		return e
	}
	e = areaScoreService(&args, d, reply)
	d.Lock.Unlock()
	return e
}

// ApplyAreaScore 申请数目
func (t *Arith) ApplyAreaScore(ctx context.Context, args InfoRequest, reply *string) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("ApplyAreaScore", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	if !d.Lock.TryLockWithTimeout(time.Second * 3) {
		e = errors.New("请求过于频繁")
		return e
	}
	e = applyForAreaScoreService(&args, d)
	d.Lock.Unlock()
	return e
}

// AgreeAreaScore 同意数目
func (t *Arith) AgreeAreaScore(ctx context.Context, args InfoRequest, reply *EndResponse) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("AgreeAreaScore", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	if !d.Lock.TryLockWithTimeout(time.Second * 3) {
		e = errors.New("请求过于频繁")
		return e
	}
	e = agreeAreaScoreService(&args, d, reply)
	d.Lock.Unlock()
	return e
}

// RejectAreaScore 拒绝数目
func (t *Arith) RejectAreaScore(ctx context.Context, args InfoRequest, reply *string) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("RejectAreaScore", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	if !d.Lock.TryLockWithTimeout(time.Second * 3) {
		e = errors.New("请求过于频繁")
		return e
	}
	e = rejectAreaScoreService(&args, d)
	d.Lock.Unlock()
	return e
}

// CallEnd 修改结果
func (t *Arith) CallEnd(ctx context.Context, args CallEndRequest, reply *string) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("CallEnd", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	e = callEndService(&args)
	return e
}

// ApplySummation 申请和棋
func (t *Arith) ApplySummation(ctx context.Context, args InfoRequest, reply *string) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("ApplySummation", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	if !d.Lock.TryLockWithTimeout(time.Second * 3) {
		e = errors.New("请求过于频繁")
		return e
	}
	e = applyForSummationService(&args, d)
	d.Lock.Unlock()
	return e
}

// AgreeSummation 同意和棋
func (t *Arith) AgreeSummation(ctx context.Context, args InfoRequest, reply *EndResponse) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("AgreeSummation", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	if !d.Lock.TryLockWithTimeout(time.Second * 3) {
		e = errors.New("请求过于频繁")
		return e
	}
	e = agreeSummationService(&args, d, reply)
	d.Lock.Unlock()
	return e
}

// RejectSummation 拒绝和棋
func (t *Arith) RejectSummation(ctx context.Context, args InfoRequest, reply *string) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("RejectSummation", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	if !d.Lock.TryLockWithTimeout(time.Second * 3) {
		e = errors.New("请求过于频繁")
		return e
	}
	e = rejectSummationService(&args, d)
	d.Lock.Unlock()
	return e
}

// CanPlay
func (t *Arith) CanPlay(ctx context.Context, args InfoRequest, reply *CanPlayResponse) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("CanPlay", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	e = canPlayService(&args, d, reply)
	return e
}

// InnerAreaScore 内部数目
func (t *Arith) InnerAreaScore(ctx context.Context, args OwnershipRequest, reply *InnerAreaScoreResponse) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("InnerAreaScore", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	b, _, err := wq.Load(args.SGF)
	if err != nil {
		e = err
		return e
	}
	board := b.Board()
	bScore, wScore, endScore, controversyCount, _, err := areaScore(args.GameId, board)
	if err != nil {
		e = err
		return e
	}
	var resp InnerAreaScoreResponse
	resp.BScore = bScore
	resp.WScore = wScore
	resp.EndScore = endScore
	resp.ControversyCount = controversyCount
	*reply = resp
	return nil
}

// Ownership 形势判断
func (t *Arith) Ownership(ctx context.Context, args OwnershipRequest, reply *OwnerShipResponse) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("Ownership", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	b, _, err := wq.Load(args.SGF)
	if err != nil {
		e = err
		return e
	}
	board := b.Board()
	anyData := board.GetKataGoAnalysisData(args.GameId)
	e = gameBus.StartAnalysisScore(anyData)
	if e != nil {
		return e
	}
	for i := 0; i <= 10; i++ {
		d, _ := redisUtils.Get(fmt.Sprintf(gameBus.SCOREREDISRESULTKEY, fmt.Sprintf("%d", args.GameId)))
		if d == nil {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		var res gameBus.AnalysisScoreResult
		e = json.Unmarshal(d, &res)
		if e != nil {
			return e
		}
		var analysisRes gameBus.AnalysisScoreData
		z, e := utils.UnzipString(res.Data)
		if e != nil {
			return e
		}
		e = json.Unmarshal([]byte(z), &analysisRes)
		if e != nil {
			return e
		}
		var resp OwnerShipResponse
		colorOwnerShip := wq.GetColorOwnerShip(analysisRes.Ownership, board.Size)
		resp.Ownership = *colorOwnerShip
		resp.AnalysisData = analysisRes
		*reply = resp
		return nil
	}
	e = errors.New("数目失败")
	return e
}

// ForceReloadSGF 强制重新加载SGF
func (t *Arith) ForceReloadSGF(ctx context.Context, args ForceReloadSgfRequest, reply *string) error {
	var e error
	defer func(startTime int64) {
		logger.Logger("ForceReloadSGF", logger.INFO, e, fmt.Sprintf("耗时:%d milliseconds", time.Now().UnixMilli()-startTime))
	}(time.Now().UnixMilli())
	d, ok := cache.CachaDataMap.Get(cache.GetKey(args.GameId))
	if !ok {
		e = errors.New("对弈信息不存在")
		return e
	}
	if !d.Lock.TryLockWithTimeout(time.Second * 3) {
		e = errors.New("请求过于频繁")
		return e
	}
	e = forceReloadService(&args)
	d.Lock.Unlock()
	return e
}
