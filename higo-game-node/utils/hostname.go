package utils

import (
	"fmt"
	"higo-game-node/config"
	"os"
)

func GetHostName() string {
	if config.RunMode == "debug" {
		return "game-node-1"
	}
	hostname, _ := os.Hostname()
	return hostname
}

var gameHost []string

func InitGameHost() {
	for i := 0; i < config.Conf.Pod.GamePodNumb; i++ {
		host := fmt.Sprintf("%s-%d", config.Conf.Pod.GamePodName, i)
		gameHost = append(gameHost, host)
	}
}

func FindPodHost(a string) string {
	hash := 0
	for i := 0; i < len(a); i++ {
		hash += int(a[i])
	}
	index := hash % len(gameHost)
	if index >= 0 && index < len(gameHost) {
		return gameHost[index]
	}
	return gameHost[0]
}
