package lock

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	lockNx   = "log_lock"
	baseTime = time.Second * 50
	outTime  = time.Second * 40
)

var Rds *redis.Client

type DistributeLock struct {
	Flag  string
	DLock *sync.Mutex
}

//getLock 用户，执行逻辑函数，带守护进程的分布式锁,成功执行反馈true
//user标识谁在使用这把锁，fc为该锁执行的逻辑，name对应user
func getLock(user string, fc func(name string)) bool {
	//锁的标识得是有一个身份辨认,锁获取成功才继续
	ok, err := Rds.SetNX(lockNx, user, baseTime).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{"setnx": err.Error()}).Error("redis")
		return false
	}
	if !ok {
		return false
	}
	out := make(chan bool)
	//这是一个守护进程，检查是否完成执行，如果到锁的时间了，还没有完成，那就重新设置一个时间（续命），完成后到时间自动退出
	go func() {
		for {
			select {
			case <-out:
				//加锁doucle check
				Rds.Del(lockNx) //成功后释放锁
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case <-time.After(outTime):
				if Rds.Get(lockNx).Val() != user { //只有自己的锁，才能给他续命//保险操作
					return
				}
				Rds.Expire(lockNx, time.Minute)
				// case _, ok := <-out:
				// 	if !ok {
				// 		return
				// 	}
			}
		}
	}()
	defer func() { //fc错误收集
		if err := recover(); err != nil {
			go func() { out <- true; close(out) }() //说明执行fc异常了，释放锁
			logrus.WithFields(logrus.Fields{"Execution function": err}).Error("DistributeLock")
		}
	}()
	fc(user)
	go func() { out <- true; close(out) }() //直接退出守护进程,这里要关闭out，防止select中两个同时到期
	//唯一不足就是上面这两步中，在刚好延时操作和out信号同时发生时，out信号被忽略，select选择执行select的第一个选项，然而，这里的close还没执行到的时候
	//这时候就发生out信号被忽略，关闭管道没执行，导致select第一步仍然续命一次，知道下一次检查超时，发现通道关闭才会释放锁（就是可能锁的占用时间多一个延时检查的时间）
	return true
}
