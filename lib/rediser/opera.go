package rediser

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

//Open 打开redis连接
func Open(addr, pwd string, dbevel int) *redis.Client {
	// redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", Password: "", DB: 0})
	rds := redis.NewClient(&redis.Options{Addr: addr, Password: pwd, DB: dbevel})
	_, err := rds.Ping().Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{"connect": err}).Error("redis")
	}
	logrus.WithFields(logrus.Fields{"connect": addr}).Info("redis")
	return rds
}

//Close 关闭
func Close() {

}

// GetOnlineMember 获取所有在线成员
func GetOnlineMember(rd *redis.Client) {

}

//GetUser 获取用户信息
func GetUser(rd *redis.Client) {
	// res, err := rd.HMGet("stb").Result()
	// if err != nil {
	// 	logrus.WithFields(logrus.Fields{"getuser": err}).Error("redisErr")
	// }
	// logrus.WithFields(logrus.Fields{"getuser": res}).Info("redis")
	na, err := rd.Get("name").Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{"getuser": err}).Error("redisErr")
	}
	logrus.Info("name:", na)
}

//SetUser 设置用户信息
func SetUser(rd *redis.Client) {

	rd.Set("name", "stb", 5)

	// if err != nil {
	// 	logrus.WithFields(logrus.Fields{"setuser": err}).Error("redisErr")
	// }
}
