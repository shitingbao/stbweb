package core

import (
	"time"

	"github.com/sirupsen/logrus"
)

const (
	lockNx   = "log_lock"
	baseTime = time.Second * 50
	outTime  = time.Second * 40
)

//DistributeLock 用户，执行逻辑函数，带守护进程的分布式锁,成功执行反馈true
func DistributeLock(user string, fc func()) bool {
	//锁的标识得是有一个身份辨认,锁获取成功才继续
	ok, err := Rds.SetNX(lockNx, "shitingbao", baseTime).Result()
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
			case <-time.After(outTime):
				if Rds.Get(lockNx).Val() != "shitingbao" { //只有自己的锁，才能给他续命//保险操作，非必须
					return
				}
				//判断out是否关闭，因为这里可能两个都符合，随机执行的时候，其实已经完成任务，但是选择了续命锁的操作
				if _, ok := <-out; !ok { //通道关闭说明已经完成
					Rds.Del(lockNx)
					return
				}
				Rds.Expire(lockNx, time.Minute)
			case <-out:
				Rds.Del(lockNx) //成功后释放锁
				return
			}
		}
	}()
	fc()
	go func() { out <- true; close(out) }() //直接退出守护进程,这里要关闭out，防止select中两个同时到期
	return true
}
