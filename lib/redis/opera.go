package redis

import "github.com/go-redis/redis"

//Open 打开redis连接
func Open(addr, pwd string, dbevel int) {
	// redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", Password: "", DB: 0})
	redis.NewClient(&redis.Options{Addr: addr, Password: pwd, DB: dbevel})
}

//Close 关闭
func Close() {

}
