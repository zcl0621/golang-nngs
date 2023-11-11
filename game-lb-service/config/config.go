package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var RunMode string = retrieveEnvOrDefault("RUN_MODE", "debug")
var Conf *configYaml

func InitConf() {
	if Conf == nil { // 避免反复读取配置文件
		// 读取配置文件
		if err := ReadConfig(); err != nil {
			log.Fatalf("读取配置文件失败: %s", err.Error())
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
	Http             httpYaml `yaml:"http"`
	GamePodNamespace string   `yaml:"game_pod_namespace"`
	GamePodName      string   `yaml:"game_pod_name"`
	GamePodService   string   `yaml:"game_pod_service"`
	GamePodNumb      int      `yaml:"game_pod_numb"`
}

type httpYaml struct {
	Port string `yaml:"port"`
}

func (c *configYaml) getConf(path string) error {
	if yamlFile, err := ioutil.ReadFile(path); err != nil {
		return err
	} else {
		return yaml.UnmarshalStrict(yamlFile, c)
	}
}
