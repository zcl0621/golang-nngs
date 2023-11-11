package config

import (
	"gopkg.in/yaml.v2"
	"higo-game-node/logger"
	"io/ioutil"
	"os"
)

var RunMode string = retrieveEnvOrDefault("RUN_MODE", "debug")
var Conf *configYaml

func InitConf() {
	if Conf == nil { // 避免反复读取配置文件
		// 读取配置文件
		if err := ReadConfig(); err != nil {
			logger.Logger("main 读取配置文件失败", "error", err, "")
			os.Exit(2)
			return
		}
	}
}

// ReadConfig 读取配置文件
func ReadConfig() error {
	var configPath string
	if RunMode == "debug" {
		configPath = "config.yaml"
	} else {
		configPath = "/etc/config.yaml"
	}
	var configData configYaml
	if err := configData.getConf(configPath); err != nil {
		return err
	} else {
		Conf = &configData
		return nil
	}
}

func retrieveEnvOrDefault(key string, defaultValue string) string {
	result := os.Getenv(key)
	if len(result) == 0 {
		result = defaultValue
	}
	return result
}

type configYaml struct {
	Http         httpYaml     `yaml:"http"`
	DataBase     databaseYaml `yaml:"database"`
	ThirdService thirdService `yaml:"third_service"`
	Redis        redisYaml    `yaml:"redis"`
	Pod          podYaml      `yaml:"pod"`
	Rule         ruleYaml     `yaml:"rule"`
	MongoDB      mongoDBYaml  `yaml:"mongodb"`
}

type mongoDBYaml struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type podYaml struct {
	GamePodName string `yaml:"game_pod_name"`
	GamePodNumb int    `yaml:"game_pod_numb"`
}

type thirdService struct {
	WsService      string `yaml:"ws_service"`
	AiAgent        string `yaml:"ai_agent"`
	GameBusService string `yaml:"game_bus_service"`
}
type httpYaml struct {
	Port string `yaml:"port"`
}

type databaseYaml struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type redisYaml struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Channel string `yaml:"channel"`
}

type ruleYaml struct {
	FirstMoveTime         int   `yaml:"first_move_time"`
	MoveTime              int   `yaml:"move_time"`
	AreaScoreTimeInterval int64 `yaml:"area_score_time_interval"`
	AreaScoreTime         int   `yaml:"area_score_time"`
	SummationTimeInterval int64 `yaml:"summation_time_interval"`
	SummationTime         int   `yaml:"summation_time"`
	SummationCount        int   `yaml:"summation_count"`
}

func (c *configYaml) getConf(path string) error {
	if yamlFile, err := ioutil.ReadFile(path); err != nil {
		return err
	} else {
		return yaml.UnmarshalStrict(yamlFile, c)
	}
}
