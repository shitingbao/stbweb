package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//Config 配置内容
type Config struct {
	Driver          string //数据库标识
	ConnectString   string //sql连接
	BaidubceAddress string //外地址API
	AccessToken     string //外地址APItoken
	AllowCORS       bool   //是否允许本地跨域
	LogLevel        string //log等级
	Port            string //监听端口
}

//ReadConfig 读取本地config,传入config地址路径，反馈配置对象
//因为这个config对象在外部使用过程中可能会被赋值，所以使用指针
func ReadConfig(filename string) *Config {
	config := &Config{
		LogLevel: "debug", //default value
	}
	bys, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panic(err)
	}
	if err := json.Unmarshal(bys, config); err != nil {
		log.Panic(err)
	}
	return config
}
