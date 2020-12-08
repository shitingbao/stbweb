package rediser

import (
	"github.com/go-redis/redis"
)

//注意这里的端口地址，指的是哨兵的地址

//OpenSential 单个哨兵节点
func OpenSential() {
	sentinel := redis.NewSentinelClient(&redis.Options{
		Addr:       ":" + "26379",
		MaxRetries: -1,
	})
	defer sentinel.Close()
}

//OpenFailoverSential 多个哨兵节点，具有故障转移
func OpenFailoverSential() {
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "master",
		SentinelAddrs: []string{":26379", ":26380"},
	})
	if err := rdb.Ping().Err(); err != nil {
		return
	}
	defer rdb.Close()
}
