package rediser

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

//Open 打开redis连接
func Open(addr, pwd string, dbevel int) *redis.Client {
	// redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", Password: "", DB: 0})
	defer func() {
		if err := recover(); err != nil {
			logrus.WithFields(logrus.Fields{"connect": err}).Panic("redis")
		}
	}()
	return redis.NewClient(&redis.Options{Addr: addr, Password: pwd, DB: dbevel})
}

//Close 关闭
func Close() {

}
