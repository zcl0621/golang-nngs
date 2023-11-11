package score

import (
	"encoding/json"
	"fmt"
	scoreApi "higo-game-bus/api/score"
	"higo-game-bus/config"
	"higo-game-bus/logger"
	"higo-game-bus/queue"
	"higo-game-bus/redisUtils"
	"higo-game-bus/request"
)

var SCOREREDISJOBKEY = "higo:game:bus:score:job"
var SCOREREDISRESULTKEY = "higo:game:bus:score:result:%s"

var SCOREQUEUE = queue.NewQueue(1024)

func StartJob() {
	for i := 0; i < config.Conf.ThirdService.GoScoreCount; i++ {
		go handlerJob()
	}
}

func handlerJob() {
	for {
		data := SCOREQUEUE.Pop()
		if data != nil {
			func(data interface{}) {
				defer func() {
					if err := recover(); err != nil {
						logger.Logger("score.HandlerJob", logger.ERROR, nil, fmt.Sprintf("e: %s", err))
					}
				}()
				res := data.(request.AnalysisScoreRequest)
				resp, err := scoreApi.AnalysisScoreApi(&res)
				if err != nil {
					logger.Logger("score.HandlerJob AnalysisScoreApi", logger.ERROR, err, fmt.Sprintf("data %s e: %s", res, err.Error()))
					SetJob(&res)
					return
				}
				d, _ := json.Marshal(resp)
				e := redisUtils.Set(fmt.Sprintf(SCOREREDISRESULTKEY, resp.Id), d, 5*60)
				if e != nil {
					logger.Logger("score.HandlerJob redisUtils.Set", logger.ERROR, e, fmt.Sprintf("data %s resp %v e: %s", data, resp, e.Error()))
					SetJob(&res)
					return
				}
			}(data)
		}
	}
}

func SetJob(analysisRequest *request.AnalysisScoreRequest) {
	SCOREQUEUE.Push(*analysisRequest)
}
