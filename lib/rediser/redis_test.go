package rediser

import (
	"log"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

func gettest(rd *redis.Client) {
	na, err := rd.Get("name").Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{"getuser": err}).Error("redisErr")
	}
	logrus.Info("name:", na)

	tm, err := rd.TTL("name").Result() //对应健的剩余时间
	if err != nil {
		logrus.WithFields(logrus.Fields{"nametime": err}).Error("redisErr")
	}
	logrus.Info("nametime:", tm)
	isn, err := rd.PExpire("name", time.Second*10).Result() //重新设置过期时间
	if err != nil {
		logrus.WithFields(logrus.Fields{"nametime": err}).Error("redisErr")
	}
	logrus.Info("name new time:", isn, "-", rd.TTL("name").Val())
	logrus.Info("hmuser:", rd.HMGet("hmuser", "age", "ip").Val()) //获取哈希对象内容，必须两个对应列,反馈的是一个slice，内容是对应列的值，这里就是对应age和ip的值18-192.168.1.39

	logrus.Info("userlist:", rd.LRange("userlist", 0, 10).Val()) //获取列表list数组内容，必须两个对应列

	rd.Expire("huser", time.Minute) //手动设置过期时间

	logrus.Info("huser:", rd.HGet("huser", "hname").Val())

	um, err := rd.HGetAll("user").Result() //获取所有该hash的对象值
	if err != nil {
		logrus.WithFields(logrus.Fields{"getuser": err}).Error("redisErr")
	}
	log.Println(um)
}

func settest(rd *redis.Client) {
	if err := rd.Set("name", "stb", time.Second*5).Err(); err != nil { //设置字符串key
		logrus.WithFields(logrus.Fields{"set": err}).Error("redisErr")
	}
	if err := rd.Set("name", "shtingbao", time.Second*8).Err(); err != nil { //重复设置会覆盖
		logrus.WithFields(logrus.Fields{"set": err}).Error("redisErr")
	}
	hmuser := make(map[string]interface{})
	hmuser["age"] = 18
	hmuser["ip"] = "192.168.1.39"
	rd.HMSet("hmuser", hmuser) //设置hash对象

	rd.LPush("userlist", "1", "2", "3") //设置列表list，添加123进入userlist

	rd.HSet("huser", "hname", "shitingbao")

	if err := rd.HSet("user", "a", "asdf").Err(); err != nil { //设置多个hash值
		logrus.WithFields(logrus.Fields{"set": err}).Error("redisErr")
	}
	if err := rd.HSet("user", "b", "qwer").Err(); err != nil {
		logrus.WithFields(logrus.Fields{"set": err}).Error("redisErr")
	}
	if err := rd.HSet("user", "c", "zxvc").Err(); err != nil {
		logrus.WithFields(logrus.Fields{"set": err}).Error("redisErr")
	}
}
