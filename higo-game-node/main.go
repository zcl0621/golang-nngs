package main

import (
	"flag"
	"fmt"
	"github.com/smallnest/rpcx/server"
	"higo-game-node/api/ws"
	"higo-game-node/config"
	"higo-game-node/database"
	"higo-game-node/handlerHistory"
	"higo-game-node/mongodb"
	"higo-game-node/play"
	"higo-game-node/redisUtils"
	"higo-game-node/steam"
	"higo-game-node/utils"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	if config.RunMode == "debug" || config.RunMode == "dev" {
		go pprofDebug()
	}
	config.InitConf()
	redisUtils.InitRedis()
	database.InitDatabase()
	utils.InitGameHost()
	mongodb.InitMongoDB()
	go steam.SubWsConnect()
	go ws.StartWsPushChan()
	go steam.StartNotifyChan()
	handlerHistory.InitQueue()
	play.ReloadGame()
	flag.Parse()
	addr := flag.String("addr", fmt.Sprintf("0.0.0.0:%s", config.Conf.Http.Port), "server address")
	s := server.NewServer(
		server.WithPool(1000, 1024),
		server.WithReadTimeout(time.Second*5),
		server.WithWriteTimeout(time.Second*5),
	)
	e := s.RegisterName("Arith", new(play.Arith), "")
	if e != nil {
		panic(e)
	}
	e = s.Serve("tcp", *addr)
	if e != nil {
		panic(e)
	}
}

func pprofDebug() {
	http.ListenAndServe("127.0.0.1:9999", nil)
}
