package main

import (
	"game-lb/api/node"
	"game-lb/config"
	"game-lb/server"
	"os"
)

func main() {
	config.InitConf()
	node.Init()
	if err := server.StartGinServer(); err != nil {
		os.Exit(2)
		return
	}
}
