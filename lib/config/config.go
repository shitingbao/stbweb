package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	AllowOrigin     string //允许跨域地址
	AccessTokenDate string //文字识别接口token的有效期，自动写入，不需要手动修改

	RedisAdree string //redis连接地址
	RedisPwd   string //redis连接密码
	Redislevel int    //redis等级
	RedisPort  string //redis端口号
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

//SaveConfig 重新保存配置进入config
func (cg *Config) SaveConfig() {
	f, err := os.Create("config.json")
	if err != nil {
		log.Panic(err)
	}
	json.NewEncoder(f).Encode(cg)
	defer f.Close()
}
