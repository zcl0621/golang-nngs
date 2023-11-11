package main

import (
	"higo-game-bus/config"
	"higo-game-bus/database"
	"higo-game-bus/logger"
	"higo-game-bus/maintain"
	"higo-game-bus/migrate"
	"higo-game-bus/pub"
	"higo-game-bus/redisUtils"
	"higo-game-bus/score"
	"higo-game-bus/server"
	"higo-game-bus/steam"
	"os"
)

func main() {
	config.InitConf()
	redisUtils.InitRedis()
	database.InitDatabase()
	migrate.MigrateModel()
	steam.SubGameNotify()
	pub.StartInitData()
	maintain.CheckMaintain()
	score.StartJob()
	if err := server.StartGinServer(); err != nil {
		logger.Logger("main 启动http服务失败", "error", err, "")
		os.Exit(2)
		return
	}
}
