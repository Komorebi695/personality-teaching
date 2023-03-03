package configs

import (
	"gopkg.in/yaml.v2"
	"os"
)

// AppConfig 服务端配置（组合全部配置模型）
type AppConfig struct {
	AppName  string `yaml:"app_name"`
	Port     string `yaml:"port"`
	Mode     string `yaml:"mode"`
	DataBase Mysql  `yaml:"data_base"`
	Redis    Redis  `yaml:"redis"`
}

// Mysql 配置
type Mysql struct {
	Drive    string `yaml:"drive"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Pwd      string `yaml:"pwd"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}

// Redis 配置
type Redis struct {
	NetWork  string `yaml:"net_work"`
	Addr     string `yaml:"addr"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	Prefix   string `yaml:"prefix"`
}

// InitConfig 初始化服务器配置
func InitConfig() *AppConfig {
	var config *AppConfig
	file, err := os.Open("./src/configs/config.yaml")
	if err != nil {
		panic(any(err.Error()))
	}
	// 解析
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(any(err.Error()))
	}
	return config
}
