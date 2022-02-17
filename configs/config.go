package configs

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type GlobalConfig struct {
	ItemId  string
	SecTime string
	Email   string
}

var Config GlobalConfig

func init() {
	configFileName := "configs/config.yaml"
	if data, err := ioutil.ReadFile(configFileName); err != nil {
		panic("初始化配置文件失败！")
	} else {
		yaml.Unmarshal(data, &Config)
	}
}
