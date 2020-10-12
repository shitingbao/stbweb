package main

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

var (
	rds *redis.Client
)

//带守护进程的分布式锁
func load() {
	rds = redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379", Password: "", DB: 0})
	go loadTest()
	time.Sleep(time.Second * 8)
	log.Println("out")

}

func loadTest() {
	if tm := rds.Get("tm").Val(); tm != "" { //检查锁状态
		return
	}
	out := make(chan bool)
	//这是一个守护进程，检查是否完成执行，如果到锁的时间了，还没有完成，那就重新设置一个时间（续命），完成后到时间自动退出
	go func() {
		//不能放在go外面，否则下面get可能获取不到这个shitingbao
		//锁的标识得是有一个身份辨认
		rds.Set("tm", "shitingbao", time.Second*4)
		for {
			select {
			case <-time.After(time.Second):
				tm, err := rds.Get("tm").Result()
				if err != nil {
					log.Println(err)
				}
				log.Println("tm:", tm)
				if tm != "shitingbao" { //只有自己的锁，才能给他续命
					log.Println("stop")
					return
				}
				rds.Expire("tm", time.Second*4)
				log.Println("延迟")
			case <-out:
				return
			}
		}
	}()
	log.Println("opera")
	time.Sleep(time.Second * 5)
	log.Println("complete")
	rds.Del("tm")               //成功后释放锁
	go func() { out <- true }() //直接退出守护进程
}
