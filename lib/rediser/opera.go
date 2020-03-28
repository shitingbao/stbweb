package rediser

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

const (
	//UserMerber 用户列表 redis的hash的key，保存所有成员的对应信息，需要定时清理
	UserMerber = "user@list"
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
//记录所有成员，形成一个map
//因为每个成员的name是用key作对应设置时间的，不使用集合，因为每个成员的过期时间都不一样不好统计
//所以获取所有在线成员，就把名称取出来，一一对照过期时间来反馈
func GetOnlineMember(rd *redis.Client) []string {
	um, err := rd.HGetAll(UserMerber).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{"getAlluser": err}).Error("redisErr")
		return []string{}
	}
	// log.Println(um)
	userList := []string{}
	for k, v := range um {
		if rd.Get(k).Val() == "" {
			continue
		}
		userList = append(userList, v)
	}
	return userList
}

//GetUser 获取用户信息,无用户为空字符串
func GetUser(rd *redis.Client, userkey string) string {
	name, err := rd.Get(userkey).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{"getuser": err}).Error("redisErr")
	}

	return name
}

//SetUser 设置用户信息
func SetUser(rd *redis.Client, userkey, username string) {

	if err := rd.Set(userkey, username, time.Minute*5).Err(); err != nil { //设置字符串key
		logrus.WithFields(logrus.Fields{"set": err}).Error("redisErr")
	}
	if err := rd.HSet(UserMerber, userkey, username).Err(); err != nil { //设置字符串key
		logrus.WithFields(logrus.Fields{"setuserlist": err}).Error("redisErr")
	}
}

//DelUser 删除一个用户，需要删除对应key和用户列表中的数据
func DelUser(rd *redis.Client, userkey string) {
	if err := rd.Del(userkey).Err(); err != nil {
		logrus.WithFields(logrus.Fields{"deluser": err}).Error("redisErr")
		return
	}
	if err := rd.HDel(UserMerber, userkey).Err(); err != nil {
		logrus.WithFields(logrus.Fields{"deluserlist": err}).Error("redisErr")
		return
	}
}
