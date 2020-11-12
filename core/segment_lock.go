//Package core 使用该锁注意，依赖于redis
//1.一个get，对应一个Freed，获取和释放需要对应，不然通道内的信号量达不到一致性可能会
//2.在使用时，应该给他设置以user为key的setnx
package core

import (
	"time"

	"github.com/sirupsen/logrus"
)

//SegmentLockPro 锁标识前缀
var SegmentLockPro = "segment_lock_pro_"

//CustomizeLock cap总长度，锁实际内容，key标识所属（user），bool代表是否已经使用，true标识已经使用，超时删除map内的对应key关系
type CustomizeLock struct {
	cap   int
	locks map[string]bool
	flag  chan bool
}

//NewLock 新建一个锁，放入需要使用的长度
//同时在flag队列里面放满标识待用
//注意这里不能关闭标识通道，因为需要复用，重新加入标识（比如超时没有使用该锁）
func NewLock(cap int) *CustomizeLock {
	flag := make(chan bool, cap)
	for i := 0; i < cap; i++ {
		flag <- true
	}
	// close(flag)
	return &CustomizeLock{
		cap:   cap,
		locks: make(map[string]bool),
		flag:  flag,
	}
}

//GetLock 从标识队列中获取一个锁,并加入使用对象，同时开始计时
//使用该锁时，应当设置setnx，给超时检查一个信号，说明已经使用
//使用过程中应该给setnx延续时间，以本身的user为key
//同理，在使用时如果setnx发生错误，说明已经超时，需要重新获取锁
func (c *CustomizeLock) GetLock(user string) bool {
	select {
	case <-c.flag:
		c.locks[user] = true
		go func() {
			tm := time.NewTicker(time.Second)
			select {
			case <-tm.C:
				if err := Rds.SetNX(SegmentLockPro+user, user, time.Second).Err(); err == nil { //这里设置成功说明使用者在规定时间段内没有使用这把锁
					delete(c.locks, user)
					c.flag <- true //记得放回标识
				}
			}
		}()
		return true
	default:
		return false
	}
}

// FreedLock 释放锁
func (c *CustomizeLock) FreedLock(user string) {
	delete(c.locks, user)
	if err := Rds.Del(user).Err(); err != nil {
		logrus.WithFields(logrus.Fields{"del lock": err}).Error("segment_lock")
	}
	c.flag <- true //记得放回标识、
}

//ResetLockOutTime 重置过期时间
func (c *CustomizeLock) ResetLockOutTime(user string) {
	Rds.Expire(SegmentLockPro+user, time.Second)
}
